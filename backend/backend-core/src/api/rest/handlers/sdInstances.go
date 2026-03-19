package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/go-chi/chi/v5"
)

func GetSDInstances(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceSDInstances, auth.OperationRead) {
		http.Error(w, "forbidden", http.StatusForbidden)
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
	if !auth.CanAccessOperation(r.Context(), auth.ResourceSDInstances, auth.OperationUpdate) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var input graphQLModel.SDInstanceUpdateInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateSDInstance(uint32(id), input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}
