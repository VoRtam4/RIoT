package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
)

func GetKPIResults(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceKPIResults, auth.OperationRead) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	result := domainLogicLayer.GetKPIFulfillmentCheckResults()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}
