/**
 * @file user.go
 * @brief REST handlery pro práci s uživatelskou konfigurací.
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
	"fmt"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/go-chi/chi/v5"
)

type disableUserRequest struct {
	Reason *string `json:"reason"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceUsers, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.LoadUsers()
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceUsers, auth.OperationRead); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.LoadUserByUID(uid)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	principal, ok := auth.PrincipalFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	result := domainLogicLayer.LoadUserByID(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceUsers, auth.OperationUpdate); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	var input graphQLModel.UserUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateUser(uid, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DisableUser(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceUsers, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	var req disableUserRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	result := domainLogicLayer.DisableUser(principal.UserID, uid, req.Reason)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func EnableUser(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceUsers, auth.OperationUpdate); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.EnableUser(uid)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func RevokeUserSessions(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceUsers, auth.OperationUpdate); principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	if err := domainLogicLayer.RevokeUserSessions(uid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetSessions(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceSessions, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.LoadUserSessions(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetSessionsByUser(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceUsers, auth.OperationRead); principal == nil {
		return
	}
	userUID := chi.URLParam(r, "uid")
	if userUID == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.LoadUserSessionsByUserUID(userUID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func RevokeSession(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceSessions, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	if err := domainLogicLayer.RevokeSession(principal.UserID, uid, allowForeign); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func RevokeOwnOtherSessions(w http.ResponseWriter, r *http.Request) {
	principal := authorizeOperation(w, r, auth.ResourceSessions, auth.OperationUpdate)
	if principal == nil {
		return
	}
	if principal.SessionID == nil {
		http.Error(w, fmt.Errorf("current session not available").Error(), http.StatusBadRequest)
		return
	}
	if err := domainLogicLayer.RevokeOwnOtherSessions(principal.UserID, *principal.SessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetUserConfig(w http.ResponseWriter, r *http.Request) {
	var principal *auth.Principal
	if principal = authorizeOperation(w, r, auth.ResourceUserConfig, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetUserConfig(principal.UserID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	var principal *auth.Principal
	if principal = authorizeOperation(w, r, auth.ResourceUserConfig, auth.OperationUpdate); principal == nil {
		return
	}
	var input graphQLModel.UserConfigInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateUserConfig(principal.UserID, input)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteUserConfig(w http.ResponseWriter, r *http.Request) {
	var principal *auth.Principal
	if principal = authorizeOperation(w, r, auth.ResourceUserConfig, auth.OperationDelete); principal == nil {
		return
	}
	err := domainLogicLayer.DeleteUserConfig(principal.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
