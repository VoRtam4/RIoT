package handlers

/*
import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func QueryTimeSeries(w http.ResponseWriter, r *http.Request) {

	var req sharedModel.TimeSeriesReadRequest

	json.NewDecoder(r.Body).Decode(&req)

	result := domainLogicLayer.ReadTimeSeries(req)

	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.GetPayload())
}
*/
