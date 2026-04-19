package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetKPIDefinitions(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.GetKPIDefinitions(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetKPIDefinition(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationRead)
	if principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	result := domainLogicLayer.GetKPIDefinition(principal.UserID, id)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetKPIDefinitionsBySDType(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationRead)
	if principal == nil {
		return
	}
	sdTypeID, ok := parseID(w, r)
	if !ok {
		return
	}
	result := domainLogicLayer.GetKPIDefinitionsBySDType(principal.UserID, sdTypeID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetKPIDefinitionsBySDInstace(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationRead)
	if principal == nil {
		return
	}
	sdInstanceID, ok := parseID(w, r)
	if !ok {
		return
	}
	result := domainLogicLayer.GetKPIDefinitionsBySDInstance(principal.UserID, sdInstanceID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func CreateKPIDefinition(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationCreate)
	if principal == nil {
		return
	}
	var input graphQLModel.KPIDefinitionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.CreateKPIDefinition(principal.UserID, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateKPIDefinition(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationUpdate)
	if principal == nil {
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
	result := domainLogicLayer.UpdateKPIDefinition(principal.UserID, id, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteKPIDefinition(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceKPIDefinitions, auth.OperationDelete)
	if principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	err := domainLogicLayer.DeleteKPIDefinition(principal.UserID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
