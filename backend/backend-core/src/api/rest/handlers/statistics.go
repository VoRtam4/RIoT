package handlers

/*
import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func StatisticsQuery(w http.ResponseWriter, r *http.Request) {

	var input graphQLModel.StatisticsInput

	json.NewDecoder(r.Body).Decode(&input)

	converted, _ := domainLogicLayer.MapStatisticsInputToReadRequestBody(&input, nil, nil)

	result := domainLogicLayer.Query(*converted)

	json.NewEncoder(w).Encode(result.Unwrap())
}
*/
