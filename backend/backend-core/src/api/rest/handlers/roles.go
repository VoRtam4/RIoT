/**
 * @file roles.go
 * @brief REST handlery pro práci s rolemi a přiřazením rolí uživatelům.
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
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.LoadRoles()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetUserRole(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	userID, ok := parseID(w, r)
	if !ok {
		return
	}
	result := domainLogicLayer.LoadUserRole(userID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetRole(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.LoadUserRole(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationUpdate); principal == nil {
		return
	}
	var req graphQLModel.AssignRoleInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := domainLogicLayer.AssignRoleToUser(req.UserID, req.RoleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
