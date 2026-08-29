/**
 * @file principal.go
 * @brief Model a autentizace principala pro session, API klíče a hostovský přístup.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
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

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type PrincipalType string

const (
	PrincipalUserSession PrincipalType = "user_session"
	PrincipalAPIKey      PrincipalType = "api_key"
	PrincipalGuest       PrincipalType = "guest"
)

type Principal struct {
	Type              PrincipalType
	UserID            uint32
	AllowedOperations map[string]bool
	APIKeyID          *uint32
	SessionID         *uint
	SessionUID        *string
	SynchronizedAt    time.Time
	ClientIP          string
}

const PrincipalContextKey = "principal"

func ContextWithPrincipal(ctx context.Context, p *Principal) context.Context {
	return context.WithValue(ctx, PrincipalContextKey, p)
}

func PrincipalFromContext(ctx context.Context) (*Principal, bool) {
	p, ok := ctx.Value(PrincipalContextKey).(*Principal)
	return p, ok
}

func AuthenticatePrincipal(w http.ResponseWriter, r *http.Request) (*Principal, error) {
	clientIP := extractClientIP(r)
	apiKeyRaw := extractAPIKey(r)
	if apiKeyRaw != "" {
		hash := sharedUtils.GenerateHexHash(apiKeyRaw)
		apiKeyResult := domainLogicLayer.LoadAPIKeyByHash(hash)
		if apiKeyResult.IsFailure() {
			return nil, apiKeyResult.GetError()
		}
		apiKeyOpt := apiKeyResult.GetPayload()
		if apiKeyOpt.IsEmpty() {
			return nil, fmt.Errorf("invalid api key")
		}
		return buildAPIKeyPrincipal(apiKeyOpt.GetPayload(), clientIP)
	}
	authResult := AuthenticateRequest(r)
	if authResult.IsFailure() {
		return nil, fmt.Errorf("not authenticated")
		//return buildGuestPrincipal(clientIP)
	}
	authPayload := authResult.GetPayload()
	if err := ApplyAuthenticationResult(w, authPayload); err != nil {
		return nil, err
	}
	userIDStr := authPayload.UserID
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid userID")
	}
	principal, err := buildUserPrincipal(uint32(userIDUint), clientIP)
	if err != nil {
		return nil, err
	}
	principal.SessionID = authPayload.SessionID
	principal.SessionUID = authPayload.SessionUID
	return principal, nil
}

func CanAccessOperation(principal *Principal, resource string, operation string) bool {
	if time.Since(principal.SynchronizedAt) > 30*time.Second {
		if err := refreshPrincipalPermissions(principal); err != nil {
			return false
		}
	}
	allowed, ok := principal.AllowedOperations[resource+"."+operation]
	return ok && allowed
}

func buildAPIKeyPrincipal(apiKey dllModel.APIKey, clientIP string) (*Principal, error) {
	if apiKey.Revoked {
		return nil, fmt.Errorf("api key revoked")
	}
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("api key expired")
	}
	if !isIPAllowed(clientIP, apiKey.IPRestrictions) {
		return nil, fmt.Errorf("ip not allowed")
	}
	if apiKey.ID.IsEmpty() {
		return nil, fmt.Errorf("api key missing id")
	}
	if apiKey.UserID == nil {
		return nil, fmt.Errorf("api key missing user id")
	}
	userResult := domainLogicLayer.LoadUserByID(*apiKey.UserID)
	if userResult.IsFailure() {
		return nil, userResult.GetError()
	}
	if userResult.GetPayload().Disabled {
		return nil, fmt.Errorf("user disabled")
	}
	if !allowAPIKeyRequest(apiKey.UID, apiKey.RateLimit) {
		return nil, fmt.Errorf("api key rate limit exceeded")
	}
	id := apiKey.ID.GetPayload()
	allowed, err := buildAllowedOperationsForAPIKey(*apiKey.UserID, apiKey.Permissions)
	if err != nil {
		return nil, err
	}
	if err := domainLogicLayer.TouchAPIKeyLastUsedAt(id); err != nil {
		return nil, err
	}
	return &Principal{
		Type:              PrincipalAPIKey,
		UserID:            *apiKey.UserID,
		APIKeyID:          &id,
		SessionID:         nil,
		SessionUID:        nil,
		AllowedOperations: allowed,
		SynchronizedAt:    time.Now(),
		ClientIP:          clientIP,
	}, nil
}

func buildUserPrincipal(userID uint32, clientIP string) (*Principal, error) {
	userResult := domainLogicLayer.LoadUserByID(userID)
	if userResult.IsFailure() {
		return nil, userResult.GetError()
	}
	if userResult.GetPayload().Disabled {
		return nil, fmt.Errorf("user disabled")
	}
	permsResult := domainLogicLayer.LoadUserRole(userID)
	if permsResult.IsFailure() {
		return nil, permsResult.GetError()
	}
	allowed := make(map[string]bool)
	for _, p := range permsResult.GetPayload().Permissions {
		allowed[p.UID] = true
	}
	return &Principal{
		Type:              PrincipalUserSession,
		UserID:            userID,
		AllowedOperations: allowed,
		APIKeyID:          nil,
		SessionID:         nil,
		SessionUID:        nil,
		SynchronizedAt:    time.Now(),
		ClientIP:          clientIP,
	}, nil
}

func buildGuestPrincipal(clientIP string) (*Principal, error) {
	return &Principal{
		Type:              PrincipalGuest,
		UserID:            0,
		AllowedOperations: RolePermissions[RoleGuest],
		APIKeyID:          nil,
		SessionID:         nil,
		SessionUID:        nil,
		SynchronizedAt:    time.Now(),
		ClientIP:          clientIP,
	}, nil
}

func refreshPrincipalPermissions(principal *Principal) error {
	switch principal.Type {

	case PrincipalUserSession:
		newPrincipal, err := buildUserPrincipal(principal.UserID, principal.ClientIP)
		if err != nil {
			return err
		}
		principal.AllowedOperations = newPrincipal.AllowedOperations
		principal.SynchronizedAt = time.Now()
		return nil

	case PrincipalAPIKey:
		if principal.APIKeyID == nil {
			return fmt.Errorf("missing api key id")
		}
		result := domainLogicLayer.LoadAPIKeyByID(principal.UserID, *principal.APIKeyID)
		if result.IsFailure() {
			return result.GetError()
		}
		apiKeyResult := result.GetPayload()
		if apiKeyResult.ExpiresAt != nil {
			end, err := time.Parse(time.RFC3339Nano, *apiKeyResult.ExpiresAt)
			if err != nil {
				return fmt.Errorf("time conversion failed")
			}
			if end.Before(time.Now()) {
				return fmt.Errorf("api key expired")
			}
		}
		if !isIPAllowed(principal.ClientIP, apiKeyResult.IPRestrictions) {
			return fmt.Errorf("ip not allowed")
		}
		allowed, err := buildAllowedOperationsForAPIKey(principal.UserID, apiKeyResult.Permissions)
		if err != nil {
			return err
		}
		principal.AllowedOperations = allowed
		principal.SynchronizedAt = time.Now()
		return nil

	case PrincipalGuest:
		return nil

	default:
		return fmt.Errorf("unsupported principal type")
	}
}

func buildAllowedOperationsForAPIKey(userID uint32, keyPermissions []string) (map[string]bool, error) {
	roleResult := domainLogicLayer.LoadUserRole(userID)
	if roleResult.IsFailure() {
		return nil, roleResult.GetError()
	}
	rolePermissions := make(map[string]bool, len(roleResult.GetPayload().Permissions))
	for _, permission := range roleResult.GetPayload().Permissions {
		rolePermissions[permission.UID] = true
	}
	allowed := make(map[string]bool, len(keyPermissions))
	for _, permission := range keyPermissions {
		if rolePermissions[permission] {
			allowed[permission] = true
		}
	}
	return allowed, nil
}
