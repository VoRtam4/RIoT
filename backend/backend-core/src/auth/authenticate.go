package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type AuthenticatedRequest struct {
	UserID            string
	NewSessionJWT     string
	NewRefreshToken   string
	RefreshTokenUntil time.Time
	Refreshed         bool
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
			return sharedUtils.NewSuccessResult(AuthenticatedRequest{
				UserID: subject.GetPayload(),
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
	})
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
