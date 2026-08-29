/**
 * @file rawDataPoint.go
 * @brief REST handlery pro čtení uložených raw datových bodů.
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
	"github.com/go-chi/chi/v5"
)

func GetRawDataPointsBySDType(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceRawData, auth.OperationRead)
	if principal == nil {
		return
	}

	uid := chi.URLParam(r, "uid")
	result := domainLogicLayer.GetRawDataPointsBySDTypeUID(uid)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetRawDataPoint(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceRawData, auth.OperationRead)
	if principal == nil {
		return
	}
	uid := r.URL.Query().Get("uid")
	result := domainLogicLayer.GetRawDataPointBySDInstanceUID(uid)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}
