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
)

func GetRawDataPointsBySDType(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceRawData, auth.OperationRead)
	if principal == nil {
		return
	}

	sdTypeID, ok := parseID(w, r)
	if !ok {
		return
	}

	result := domainLogicLayer.GetRawDataPointsBySDType(sdTypeID)
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
	sdInstanceID, ok := parseID(w, r)
	if !ok {
		http.Error(w, "invalid sdInstanceID", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.GetRawDataPoint(sdInstanceID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}
