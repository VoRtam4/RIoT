package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetSDInstanceGroups(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceSDInstances, auth.OperationRead) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	result := domainLogicLayer.GetSDInstanceGroups()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetSDInstanceGroup(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceSDInstances, auth.OperationRead) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	result := domainLogicLayer.GetSDInstanceGroup(uint32(id))
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func CreateSDInstanceGroup(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceSDInstances, auth.OperationUpdate) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var input graphQLModel.SDInstanceGroupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.CreateSDInstanceGroup(input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateSDInstanceGroup(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceSDInstances, auth.OperationUpdate) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	var input graphQLModel.SDInstanceGroupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateSDInstanceGroup(uint32(id), input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteSDInstanceGroup(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceSDInstances, auth.OperationUpdate) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	err := domainLogicLayer.DeleteSDInstanceGroup(uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
