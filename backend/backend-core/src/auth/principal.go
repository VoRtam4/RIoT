package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
)

type PrincipalType string

const (
	PrincipalUserSession PrincipalType = "user_session"
	PrincipalAPIKey      PrincipalType = "api_key"
)

type Principal struct {
	Type              PrincipalType
	UserID            uint32
	AllowedOperations map[string]bool
	APIKeyID          *uint32
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
		hash := hashAPIKey(apiKeyRaw)
		db := dbClient.GetRelationalDatabaseClientInstance()
		apiKeyResult := db.LoadAPIKeyByHash(hash)
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
		return nil, authResult.GetError()
	}
	if err := ApplyAuthenticationResult(w, authResult.GetPayload()); err != nil {
		return nil, err
	}
	userIDStr := authResult.GetPayload().UserID
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid userID")
	}
	return buildUserPrincipal(uint32(userIDUint), clientIP)
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

func buildAPIKeyPrincipal(apiKey dbModel.APIKeyEntity, clientIP string) (*Principal, error) {
	if apiKey.Revoked {
		return nil, fmt.Errorf("api key revoked")
	}
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("api key expired")
	}
	if !isIPAllowed(clientIP, apiKey.IPRestrictions) {
		return nil, fmt.Errorf("ip not allowed")
	}
	allowed := make(map[string]bool)
	for _, p := range apiKey.Role.Permissions {
		allowed[p.Label] = true
	}
	return &Principal{
		Type:              PrincipalAPIKey,
		UserID:            apiKey.UserID,
		APIKeyID:          &apiKey.ID,
		AllowedOperations: allowed,
		SynchronizedAt:    time.Now(),
		ClientIP:          clientIP,
	}, nil
}

func buildUserPrincipal(userID uint32, clientIP string) (*Principal, error) {
	permsResult := domainLogicLayer.LoadPermissionsForUser(userID)
	if permsResult.IsFailure() {
		return nil, permsResult.GetError()
	}
	allowed := make(map[string]bool)
	for _, p := range permsResult.GetPayload() {
		allowed[p.Label] = true
	}
	return &Principal{
		Type:              PrincipalUserSession,
		UserID:            userID,
		AllowedOperations: allowed,
		APIKeyID:          nil,
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
		apiKeyResult := domainLogicLayer.LoadAPIKeyByID(*principal.APIKeyID)
		if apiKeyResult.IsFailure() {
			return apiKeyResult.GetError()
		}
		apiKeyOpt := apiKeyResult.GetPayload()
		if apiKeyOpt.IsEmpty() {
			return fmt.Errorf("api key not found")
		}
		apiKey := apiKeyOpt.GetPayload()
		if apiKey.Revoked {
			return fmt.Errorf("api key revoked")
		}
		if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
			return fmt.Errorf("api key expired")
		}
		if !isIPAllowed(principal.ClientIP, apiKey.IPRestrictions) {
			return fmt.Errorf("ip not allowed")
		}
		allowed := make(map[string]bool)
		for _, p := range apiKey.Role.Permissions {
			allowed[p.Label] = true
		}
		principal.UserID = apiKey.UserID
		principal.AllowedOperations = allowed
		principal.SynchronizedAt = time.Now()
		return nil

	default:
		return fmt.Errorf("unsupported principal type")
	}
}
