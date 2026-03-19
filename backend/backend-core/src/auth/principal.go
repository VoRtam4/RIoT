package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
)

type PrincipalType string

const (
	PrincipalUserSession PrincipalType = "user_session"
	PrincipalAPIKey      PrincipalType = "api_key"
)

type Principal struct {
	Type              PrincipalType
	UserID            uint32
	Role              string
	AllowedOperations map[string]bool
	APIKeyID          *uint32
}

const PrincipalContextKey = "principal"

func ContextWithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, PrincipalContextKey, p)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(PrincipalContextKey).(Principal)
	return p, ok
}

func AuthenticatePrincipal(w http.ResponseWriter, r *http.Request) (*Principal, error) {
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
	db := dbClient.GetRelationalDatabaseClientInstance()
	roleResult := db.LoadUserRole(uint32(userIDUint))
	role := RoleUser
	if roleResult.IsFailure() {
	} else {
		roleOpt := roleResult.GetPayload()
		if roleOpt.IsPresent() {
			role = roleOpt.GetPayload()
		}
	}
	principal := &Principal{
		UserID: uint32(userIDUint),
		Role:   role,
	}
	return principal, nil
}

func CanAccessOperation(ctx context.Context, resource string, operation string) bool {
	principal, ok := PrincipalFromContext(ctx)
	if !ok {
		return false
	}
	key := resource + "." + operation
	permsResult := domainLogicLayer.LoadPermissionsForUser(principal.UserID)
	if permsResult.IsFailure() {
		return false
	}
	for _, p := range permsResult.GetPayload() {
		if p.Label == key {
			return true
		}
	}
	return false
}
