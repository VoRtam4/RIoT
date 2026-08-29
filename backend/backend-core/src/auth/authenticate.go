/**
 * @file authenticate.go
 * @brief Autentizace HTTP požadavku pomocí session a refresh tokenů.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type AuthenticatedRequest struct {
	UserID            string
	NewSessionJWT     string
	NewRefreshToken   string
	RefreshTokenUntil time.Time
	Refreshed         bool
	SessionID         *uint
	SessionUID        *string
}

func AuthenticateRequest(r *http.Request) sharedUtils.Result[AuthenticatedRequest] {
	if !jwtAuthenticationMiddlewareEnabled {
		return sharedUtils.NewSuccessResult(AuthenticatedRequest{
			UserID: "1",
		})
	}
	if isCookieSet(r, SessionJWTCookieIdentifier) {
		sessionJWTString := getSessionJWTCookieValue(r).GetPayload()
		sessionJWT, err := ParseJWT(sessionJWTString)
		if err != nil {
			return sharedUtils.NewFailureResult[AuthenticatedRequest](fmt.Errorf("failed to parse JWT: %w", err))
		}
		if IsJWTValid(sessionJWT) {
			subject := extractSubjectFromJWT(sessionJWTString)
			if subject.IsFailure() {
				return sharedUtils.NewFailureResult[AuthenticatedRequest](subject.GetError())
			}
			sessionID, sessionUID := resolveCurrentSessionFromRequest(r, subject.GetPayload())
			return sharedUtils.NewSuccessResult(AuthenticatedRequest{
				UserID:     subject.GetPayload(),
				SessionID:  sessionID,
				SessionUID: sessionUID,
			})
		}
		timeUntilExpiry := getTimeUntilJWTExpiry(sessionJWT)
		if timeUntilExpiry.IsFailure() {
			return sharedUtils.NewFailureResult[AuthenticatedRequest](timeUntilExpiry.GetError())
		}
		if timeUntilExpiry.GetPayload() > 0 {
			return sharedUtils.NewFailureResult[AuthenticatedRequest](fmt.Errorf("invalid JWT"))
		}
	}
	refreshToken := getRefreshTokenCookieValue(r).GetPayloadOrDefault("")
	if refreshToken == "" {
		return sharedUtils.NewFailureResult[AuthenticatedRequest](fmt.Errorf("missing refresh token"))
	}
	refreshTokenHash := sharedUtils.GenerateHexHash(refreshToken)
	raw, err, _ := sameOriginExpiredSessionJWTRequestGroup.Do(refreshTokenHash, func() (any, error) {
		return performSessionRefresh(refreshTokenHash)
	})
	if err != nil {
		return sharedUtils.NewFailureResult[AuthenticatedRequest](err)
	}
	res := raw.(*sessionRefreshResult)
	subject := extractSubjectFromJWT(res.newSessionJWT)
	if subject.IsFailure() {
		return sharedUtils.NewFailureResult[AuthenticatedRequest](subject.GetError())
	}
	return sharedUtils.NewSuccessResult(AuthenticatedRequest{
		UserID:            subject.GetPayload(),
		NewSessionJWT:     res.newSessionJWT,
		NewRefreshToken:   res.newRefreshToken,
		RefreshTokenUntil: res.refreshTokenExpiresAt,
		Refreshed:         true,
		SessionID:         res.sessionID,
		SessionUID:        res.sessionUID,
	})
}

func resolveCurrentSessionFromRequest(r *http.Request, userID string) (*uint, *string) {
	refreshToken := getRefreshTokenCookieValue(r).GetPayloadOrDefault("")
	if refreshToken == "" {
		return nil, nil
	}
	userIDUint, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		return nil, nil
	}
	refreshTokenHash := sharedUtils.GenerateHexHash(refreshToken)
	result := dbClient.GetRelationalDatabaseClientInstance().LoadUserSessionBasedOnRefreshTokenHash(refreshTokenHash)
	if result.IsFailure() || result.GetPayload().IsEmpty() {
		return nil, nil
	}
	return sessionIdentityFromDLL(result.GetPayload().GetPayload(), uint(userIDUint))
}

func sessionIdentityFromDLL(session dllModel.UserSession, userID uint) (*uint, *string) {
	if session.UserID != userID || session.Revoked || time.Until(session.ExpiresAt) <= 0 || session.ID.IsEmpty() {
		return nil, nil
	}
	sessionID := session.ID.GetPayload()
	sessionUID := session.UID
	return &sessionID, &sessionUID
}

func ApplyAuthenticationResult(w http.ResponseWriter, authResult AuthenticatedRequest) error {
	if !authResult.Refreshed {
		return nil
	}
	if err := setupSessionJWTCookie(w, authResult.NewSessionJWT); err != nil {
		return err
	}
	setupRefreshTokenCookie(w, authResult.NewRefreshToken, time.Until(authResult.RefreshTokenUntil))
	return nil
}

func ContextWithUser(ctx context.Context, authResult AuthenticatedRequest) context.Context {
	return context.WithValue(ctx, UserIdContextIdentifier, authResult.UserID)
}
