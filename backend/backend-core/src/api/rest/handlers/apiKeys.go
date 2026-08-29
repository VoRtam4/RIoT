/**
 * @file apiKeys.go
 * @brief REST handlery pro správu API klíčů.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

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
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationRead)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.LoadAPIKeyByUID(principal.UserID, uid)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetAPIKeysByUser(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceUsers, auth.OperationRead); principal == nil {
		return
	}
	if principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationRead); principal == nil {
		return
	}
	userUID := chi.URLParam(r, "uid")
	if userUID == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.LoadAPIKeysByUserUID(userUID)
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
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	var input graphQLModel.APIKeyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := domainLogicLayer.UpdateAPIKeyForUser(principal.UserID, uid, input)
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

func RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.RevokeAPIKey(principal.UserID, uid, allowForeign)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func RotateAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.RotateAPIKey(principal.UserID, uid, allowForeign)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"key": result.GetPayload(),
	})
}

func UpdateAPIKeyPermissions(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	var req struct {
		PermissionUIDs []string `json:"permissionUIDs"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.UpdateAPIKeyPermissions(principal.UserID, uid, req.PermissionUIDs, allowForeign)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateAPIKeyRestrictions(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	var input graphQLModel.APIKeyRestrictionsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.UpdateAPIKeyRestrictions(principal.UserID, uid, input, allowForeign)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationDelete)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	err := domainLogicLayer.DeleteAPIKeyForUser(principal.UserID, uid)
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
