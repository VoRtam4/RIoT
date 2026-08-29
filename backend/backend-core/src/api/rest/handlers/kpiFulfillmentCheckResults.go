/**
 * @file kpiFulfillmentCheckResults.go
 * @brief REST handlery pro práci s výsledky vyhodnocení KPI.
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
	"github.com/go-chi/chi/v5"
)

func GetKPIResults(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIResults, auth.OperationRead)
	if principal == nil {
		return
	}

	result := domainLogicLayer.GetKPIFulfillmentCheckResults(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetKPIResultsByKPI(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIResults, auth.OperationRead)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	result := domainLogicLayer.GetKPIFulfillmentCheckResultsByKPIUID(principal.UserID, uid)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetKPIResult(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIResults, auth.OperationRead)
	if principal == nil {
		return
	}

	var input graphQLModel.KPIFulfillmentCheckResultRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result := domainLogicLayer.GetKPIFulfillmentCheckResult(principal.UserID, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.GetPayload())
}
