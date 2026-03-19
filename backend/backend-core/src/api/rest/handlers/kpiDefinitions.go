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

func GetKPIDefinitions(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceKPIDefinitions, auth.OperationRead) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	result := domainLogicLayer.GetKPIDefinitions()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetKPIDefinition(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceKPIDefinitions, auth.OperationRead) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.GetKPIDefinition(uint32(id))
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func CreateKPIDefinition(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceKPIDefinitions, auth.OperationCreate) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var input graphQLModel.KPIDefinitionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.CreateKPIDefinition(input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateKPIDefinition(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceKPIDefinitions, auth.OperationUpdate) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var input graphQLModel.KPIDefinitionInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateKPIDefinition(uint32(id), input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteKPIDefinition(w http.ResponseWriter, r *http.Request) {
	if !auth.CanAccessOperation(r.Context(), auth.ResourceKPIDefinitions, auth.OperationDelete) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err = domainLogicLayer.DeleteKPIDefinition(uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
