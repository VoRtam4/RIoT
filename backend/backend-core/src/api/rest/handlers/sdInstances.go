package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetSDInstances(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceSDInstances, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetSDInstances()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateSDInstance(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceSDInstances, auth.OperationUpdate); principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input graphQLModel.SDInstanceUpdateInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateSDInstance(id, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}
