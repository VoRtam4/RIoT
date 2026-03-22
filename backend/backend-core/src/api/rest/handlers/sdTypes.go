package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetSDTypes(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceSDTypes, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetSDTypes()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetSDType(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceSDTypes, auth.OperationRead); principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	result := domainLogicLayer.GetSDType(id)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func CreateSDType(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceSDTypes, auth.OperationCreate); principal == nil {
		return
	}
	var input graphQLModel.SDTypeInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.CreateSDType(input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteSDType(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceSDTypes, auth.OperationDelete); principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	err := domainLogicLayer.DeleteSDType(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
