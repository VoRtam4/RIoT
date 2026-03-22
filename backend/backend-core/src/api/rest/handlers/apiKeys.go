package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetAPIKeys(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.LoadAPIKeysForUser(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(result.GetPayload())
}

func GetAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationCreate)
	if principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	result := domainLogicLayer.LoadAPIKeyByID(principal.UserID, id)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationCreate)
	if principal == nil {
		return
	}
	var input graphQLModel.APIKeyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.CreateAPIKey(principal.UserID, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"key": result.GetPayload(),
	})
}

func UpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var input graphQLModel.APIKeyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := domainLogicLayer.UpdateAPIKeyForUser(principal.UserID, id, input)
	if err != nil {
		switch err.Error() {
		case "not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "forbidden":
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationDelete)
	if principal == nil {
		return
	}
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	err := domainLogicLayer.DeleteAPIKeyForUser(principal.UserID, id)
	if err != nil {
		switch err.Error() {
		case "not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "forbidden":
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
