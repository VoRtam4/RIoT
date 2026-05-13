/**
 * @file helpers.go
 * @brief Pomocné funkce REST vrstvy pro autorizaci a parsování vstupních hodnot.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package handlers

import (
	"net/http"
	"strconv"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/go-chi/chi/v5"
)

func authorizeOperation(w http.ResponseWriter, r *http.Request, operation string, opType string) *auth.Principal {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return nil
	}
	if !auth.CanAccessOperation(principal, operation, opType) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return nil
	}
	return principal
}

func parseID(w http.ResponseWriter, r *http.Request) (uint32, bool) {
	idParam := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return 0, false
	}
	return uint32(id64), true
}
