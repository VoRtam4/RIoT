/**
 * @file helpers.go
 * @brief Pomocná logika pro OAuth2 přihlášení, tvorbu uživatelů a přiřazení výchozích rolí.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní pomocná autentizační logika pro OAuth2 tok a uživatelské záznamy.
 * - Vojtěch Hubáček: doplnění práce s rolemi a přiřazování rolí uživatelům během autentizace.
 *
 * @ingroup riot_backend_core
 */
package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"google.golang.org/api/idtoken"
)

var (
	allowedOrigins        = sharedUtils.NewSetFromSlice(strings.Split(sharedUtils.GetEnvironmentVariableValue("ALLOWED_ORIGINS").GetPayloadOrDefault("http://localhost:8080,http://localhost:1234"), ","))
	rootAdminEmail        = sharedUtils.GetEnvironmentVariableValue("ROOT_ADMIN_EMAIL").GetPayloadOrDefault("")
	AdminRoleID    uint32 = 0
	UserRoleID     uint32 = 0
)

// ----- types -----

type oauth2OIDCFlowState struct {
	RandomState string `json:"randomState"`
	RedirectUrl string `json:"redirectUrl"`
}

type idTokenData struct {
	oauth2ProviderIssuedID string
	email                  string
	name                   sharedUtils.Optional[string]
	profileImageURL        sharedUtils.Optional[string]
}

// ----- functions -----

func handleRedirectUrl(r *http.Request) (string, error) {
	query := r.URL.Query()
	rawRedirectUrl := query.Get("redirect")
	if rawRedirectUrl == "" {
		return "", errors.New("missing redirect url (?redirect=...)")
	}
	decodedRedirectUrl, err := url.QueryUnescape(rawRedirectUrl)
	if err != nil {
		return "", fmt.Errorf("failed to decode the redirect url (?redirect=...): %w", err)
	}
	redirectUrl, err := url.Parse(decodedRedirectUrl)
	if err != nil {
		return "", fmt.Errorf("invalid redirect url (?redirect=...): %w", err)
	}
	if redirectUrl.Scheme != "http" && redirectUrl.Scheme != "https" {
		return "", errors.New("invalid redirect url (?redirect=...): url scheme must be 'http' or 'https'")

	}
	if redirectUrl.Host == "" {
		return "", errors.New("invalid redirect url (?redirect=...): missing host")
	}
	if !allowedOrigins.Contains(fmt.Sprintf("%s://%s", redirectUrl.Scheme, redirectUrl.Host)) {
		return "", errors.New("redirect url (?redirect=...) is not among allowed origins")
	}
	return decodedRedirectUrl, nil
}

func handleUserRecordUpsert(userData idTokenData, newRefreshToken string) sharedUtils.Result[dllModel.User] {
	dbClientInstance := dbClient.GetRelationalDatabaseClientInstance()
	userLoadResult := dbClientInstance.LoadUserBasedOnOAuth2ProviderIssuedID(userData.oauth2ProviderIssuedID)
	if userLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.User](fmt.Errorf("user record upsert failure - failed to load user record: %s", userLoadResult.GetError().Error()))
	}
	role, err := resolveRoleID(userData.email)
	if err != nil {
		return sharedUtils.NewFailureResult[dllModel.User](fmt.Errorf("user record upsert failure - failed to load role: %s", err.Error()))
	}
	user := userLoadResult.GetPayload().GetPayloadOrDefault(dllModel.User{
		ID:                     sharedUtils.NewEmptyOptional[uint](),
		RoleID:                 role,
		Username:               fmt.Sprintf("google-user-%s", userData.oauth2ProviderIssuedID),
		OAuth2Provider:         sharedUtils.NewOptionalOf("google"),
		OAuth2ProviderIssuedID: sharedUtils.NewOptionalOf(userData.oauth2ProviderIssuedID),
	})
	if user.Disabled {
		return sharedUtils.NewFailureResult[dllModel.User](fmt.Errorf("user disabled"))
	}
	user.Email = userData.email
	user.Name = userData.name
	user.ProfileImageURL = userData.profileImageURL
	user.LastLoginAt = sharedUtils.NewOptionalOf(time.Now())
	persistResult := dbClientInstance.PersistUser(user)
	if persistResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.User](fmt.Errorf("user record upsert failure - failed to persist user record: %s", persistResult.GetError().Error()))
	}
	session := dllModel.UserSession{
		ID:               sharedUtils.NewEmptyOptional[uint](),
		UserID:           persistResult.GetPayload(),
		RefreshTokenHash: sharedUtils.GenerateHexHash(newRefreshToken),
		ExpiresAt:        time.Now().Add(time.Hour * 24 * 30),
		Revoked:          false,
	}
	sessionResult := dbClientInstance.PersistUserSession(session)
	if sessionResult.IsFailure() {
		return sharedUtils.NewFailureResult[dllModel.User](fmt.Errorf("failed to persist user session: %s", sessionResult.GetError().Error()))
	}
	reload := dbClientInstance.LoadUserBasedOnOAuth2ProviderIssuedID(userData.oauth2ProviderIssuedID)
	if reload.IsFailure() || reload.GetPayload().IsEmpty() {
		return sharedUtils.NewFailureResult[dllModel.User](fmt.Errorf("failed to reload user after insert"))
	}
	reloadedUser := reload.GetPayload().GetPayload()
	reloadedUser.Sessions = []dllModel.UserSession{session}
	return sharedUtils.NewSuccessResult(reloadedUser)
}

func extractIDTokenData(idTokenPayload *idtoken.Payload) sharedUtils.Result[idTokenData] {
	oauth2ProviderIssuedID := idTokenPayload.Subject
	if oauth2ProviderIssuedID == "" {
		return sharedUtils.NewFailureResult[idTokenData](errors.New("oauth2 provider (Google) issued account id not found within the id token"))
	}
	email, ok := idTokenPayload.Claims["email"].(string)
	if !ok || email == "" {
		return sharedUtils.NewFailureResult[idTokenData](errors.New("user's email address not found within the id token"))
	}
	name, _ := idTokenPayload.Claims["name"].(string)
	profileImageURL, _ := idTokenPayload.Claims["picture"].(string)
	return sharedUtils.NewSuccessResult(idTokenData{
		oauth2ProviderIssuedID: oauth2ProviderIssuedID,
		email:                  email,
		name:                   sharedUtils.Ternary[sharedUtils.Optional[string]](name != "", sharedUtils.NewOptionalOf(name), sharedUtils.NewEmptyOptional[string]()),
		profileImageURL:        sharedUtils.Ternary[sharedUtils.Optional[string]](profileImageURL != "", sharedUtils.NewOptionalOf(profileImageURL), sharedUtils.NewEmptyOptional[string]()),
	})
}

func generateRefreshToken() string {
	return sharedUtils.GenerateRandomAlphanumericString(16)
}

func resolveRoleID(email string) (uint32, error) {
	if AdminRoleID == 0 {
		adminRoleResult := domainLogicLayer.LoadRoleIDByUID(RoleAdminUID)
		if adminRoleResult.IsFailure() {
			return 0, adminRoleResult.GetError()
		}
		AdminRoleID = adminRoleResult.GetPayload()
	}
	if UserRoleID == 0 {
		userRoleResult := domainLogicLayer.LoadRoleIDByUID(RoleUserUID)
		if userRoleResult.IsFailure() {
			return 0, userRoleResult.GetError()
		}
		UserRoleID = userRoleResult.GetPayload()
	}
	if email == rootAdminEmail {
		return AdminRoleID, nil
	}
	return UserRoleID, nil
}

func hasValidSession(r *http.Request) bool {
	if !isCookieSet(r, SessionJWTCookieIdentifier) {
		return false
	}
	sessionJWTString := getSessionJWTCookieValue(r).GetPayload()
	sessionJWT, err := ParseJWT(sessionJWTString)
	if err != nil {
		return false
	}
	return IsJWTValid(sessionJWT)
}
