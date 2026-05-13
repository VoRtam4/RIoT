/**
 * @file user.go
 * @brief REST handlery pro práci s uživatelskou konfigurací.
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

func GetUserConfig(w http.ResponseWriter, r *http.Request) {
	var principal *auth.Principal
	if principal = authorizeOperation(w, r, auth.ResourceUserConfig, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetUserConfig(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	var principal *auth.Principal
	if principal = authorizeOperation(w, r, auth.ResourceUserConfig, auth.OperationUpdate); principal == nil {
		return
	}
	var input graphQLModel.UserConfigInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateUserConfig(principal.UserID, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteUserConfig(w http.ResponseWriter, r *http.Request) {
	var principal *auth.Principal
	if principal = authorizeOperation(w, r, auth.ResourceUserConfig, auth.OperationDelete); principal == nil {
		return
	}
	err := domainLogicLayer.DeleteUserConfig(principal.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
