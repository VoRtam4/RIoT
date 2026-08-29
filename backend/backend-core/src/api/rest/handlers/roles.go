/**
 * @file roles.go
 * @brief REST handlery pro práci s rolemi a přiřazením rolí uživatelům.
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

type updatePermissionLabelRequest struct {
	Label string `json:"label"`
}

func GetRoles(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.LoadRoles()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetRoleByUID(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.LoadRoleByUID(uid)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.LoadPermissions()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetUserRole(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.LoadUserRoleByUID(uid)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetRole(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.LoadUserRole(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func CreateRole(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationCreate); principal == nil {
		return
	}
	var input graphQLModel.RoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.CreateRole(input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationUpdate); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	var input graphQLModel.RoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateRole(uid, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationDelete); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	if err := domainLogicLayer.DeleteRole(uid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func CloneRole(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationCreate); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	var req struct {
		Label string `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.CloneRole(uid, req.Label)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdatePermissionLabel(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationUpdate); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	var req updatePermissionLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdatePermissionLabel(uid, req.Label)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceRoles, auth.OperationUpdate); principal == nil {
		return
	}
	var req graphQLModel.AssignRoleInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if pathUID := chi.URLParam(r, "uid"); pathUID != "" {
		req.UserUID = pathUID
	}
	err := domainLogicLayer.AssignRoleToUser(req.UserUID, req.RoleUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
