/**
 * @file middleware.go
 * @brief HTTP middleware pro autentizaci požadavků a uložení principala do kontextu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ middleware pro session autentizaci.
 * - Vojtěch Hubáček: integrace principala do kontextu požadavku a navázání middleware na novou autorizační vrstvu.
 *
 * @ingroup riot_backend_core
 */
package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"golang.org/x/sync/singleflight"
)

const (
	UserIdContextIdentifier = "userId"
)

var (
	jwtAuthenticationMiddlewareEnabled      = sharedUtils.GetFlagEnvironmentVariableValue("JWT_AUTHENTICATION_MIDDLEWARE_ENABLED").GetPayloadOrDefault(true)
	sameOriginExpiredSessionJWTRequestGroup singleflight.Group
)

func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, err := AuthenticatePrincipal(w, r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := ContextWithPrincipal(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractSubjectFromJWT(jwt string) sharedUtils.Result[string] {
	sessionJWT, err := ParseJWT(jwt)
	if err != nil {
		return sharedUtils.NewFailureResult[string](err)
	}
	subject, err := sessionJWT.Claims.GetSubject()
	if err != nil {
		return sharedUtils.NewFailureResult[string](err)
	}
	return sharedUtils.NewSuccessResult(subject)
}

type sessionRefreshResult struct {
	newSessionJWT         string
	newRefreshToken       string
	refreshTokenExpiresAt time.Time
	sessionID             *uint
	sessionUID            *string
}

func performSessionRefresh(refreshTokenHash string) (*sessionRefreshResult, error) {
	dbClientInstance := dbClient.GetRelationalDatabaseClientInstance()
	userSessionLoadResult := dbClientInstance.LoadUserSessionBasedOnRefreshTokenHash(refreshTokenHash)
	if userSessionLoadResult.IsFailure() {
		return nil, fmt.Errorf("database operation error - failed to load user session record: %w", userSessionLoadResult.GetError())
	}
	userSessionOptional := userSessionLoadResult.GetPayload()
	if userSessionOptional.IsEmpty() {
		return nil, errors.New("no user session record found based on refresh token hash")
	}
	userSession := userSessionOptional.GetPayload()
	if userSession.Revoked {
		return nil, errors.New("the refresh token has been revoked")
	}
	if time.Until(userSession.ExpiresAt) <= 0 {
		return nil, errors.New("the session has expired")
	}
	newSessionJWT, err := createSessionJWT(fmt.Sprintf("%d", userSession.UserID))
	if err != nil {
		return nil, err
	}
	newRefreshToken := sharedUtils.GenerateRandomAlphanumericString(16)
	userSession.RefreshTokenHash = sharedUtils.GenerateHexHash(newRefreshToken)
	userSessionPersistResult := dbClientInstance.PersistUserSession(userSession)
	if userSessionPersistResult.IsFailure() {
		return nil, fmt.Errorf("database operation error - failed to persist user session record: %w", userSessionPersistResult.GetError())
	}
	reloadResult := dbClientInstance.LoadUserSessionBasedOnRefreshTokenHash(userSession.RefreshTokenHash)
	if reloadResult.IsFailure() || reloadResult.GetPayload().IsEmpty() {
		return nil, errors.New("failed to reload refreshed session")
	}
	sessionID, sessionUID := sessionIdentityFromDLL(reloadResult.GetPayload().GetPayload(), userSession.UserID)
	return &sessionRefreshResult{
		newSessionJWT:         newSessionJWT,
		newRefreshToken:       newRefreshToken,
		refreshTokenExpiresAt: userSession.ExpiresAt,
		sessionID:             sessionID,
		sessionUID:            sessionUID,
	}, nil
}
