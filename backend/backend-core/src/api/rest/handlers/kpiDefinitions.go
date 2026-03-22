package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetKPIDefinitions(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationRead); principal == nil {
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
	if principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationRead); principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	result := domainLogicLayer.GetKPIDefinition(id)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func CreateKPIDefinition(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationCreate); principal == nil {
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
	if principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationUpdate); principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input graphQLModel.KPIDefinitionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateKPIDefinition(id, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteKPIDefinition(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationDelete); principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	err := domainLogicLayer.DeleteKPIDefinition(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
