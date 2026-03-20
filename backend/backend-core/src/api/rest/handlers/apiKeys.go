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
	if err := json.NewEncoder(w).Encode(result.GetPayload()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationCreate)
	if principal == nil {
		return
	}
	var input graphQLModel.CreateAPIKeyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	roleID, err := strconv.ParseUint(input.RoleID, 10, 32)
	if err != nil {
		http.Error(w, "invalid roleId", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.CreateAPIKey(
		principal.UserID,
		uint32(roleID),
		input.Label,
		input.ExpiresAt,
	)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"key": result.GetPayload(),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	idParam := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	id := uint32(id64)
	apiKeyResult := domainLogicLayer.LoadAPIKeyByID(id)
	if apiKeyResult.IsFailure() {
		http.Error(w, apiKeyResult.GetError().Error(), http.StatusInternalServerError)
		return
	}
	if apiKeyResult.GetPayload().IsEmpty() {
		http.Error(w, "api key not found", http.StatusNotFound)
		return
	}
	apiKey := apiKeyResult.GetPayload().GetPayload()
	if apiKey.UserID != principal.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var input graphQLModel.UpdateAPIKeyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if input.Label != nil {
		apiKey.Label = *input.Label
	}
	if input.RoleID != nil {
		roleID, err := strconv.ParseUint(*input.RoleID, 10, 32)
		if err != nil {
			http.Error(w, "invalid roleId", http.StatusBadRequest)
			return
		}
		apiKey.RoleID = uint32(roleID)
	}
	if input.ExpiresAt != nil {
		apiKey.ExpiresAt = input.ExpiresAt
	}
	if input.Revoked != nil {
		apiKey.Revoked = *input.Revoked
	}
	if input.RateLimit != nil {
		apiKey.RateLimit = input.RateLimit
	}
	if err := domainLogicLayer.UpdateAPIKey(apiKey); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationDelete)
	if principal == nil {
		return
	}
	idParam := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	id := uint32(id64)
	apiKeyResult := domainLogicLayer.LoadAPIKeyByID(id)
	if apiKeyResult.IsFailure() {
		http.Error(w, apiKeyResult.GetError().Error(), http.StatusInternalServerError)
		return
	}
	if apiKeyResult.GetPayload().IsEmpty() {
		http.Error(w, "api key not found", http.StatusNotFound)
		return
	}
	apiKey := apiKeyResult.GetPayload().GetPayload()
	if apiKey.UserID != principal.UserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := domainLogicLayer.DeleteAPIKey(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
