package handlers

import (
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
)

func authorizeOperation(w http.ResponseWriter, r *http.Request, operation string, opType string) *auth.Principal {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil
	}
	if !auth.CanAccessOperation(principal, auth.ResourceUserConfig, auth.OperationRead) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return nil
	}
	return principal
}
