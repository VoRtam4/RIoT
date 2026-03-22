package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationRead); principal == nil {
		return
	}
	labels := auth.GetAllRoleLabels()
	result := domainLogicLayer.LoadRolesByLabels(labels)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func GetUserRole(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationRead); principal == nil {
		return
	}
	userID, ok := parseID(w, r)
	if !ok {
		return
	}
	result := domainLogicLayer.LoadUserRole(userID)
	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result.GetPayload())
}

func AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	if principal := authorizeOperation(w, r, auth.ResourceAPIKeys, auth.OperationUpdate); principal == nil {
		return
	}
	var req graphQLModel.AssignRoleInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	err := domainLogicLayer.AssignRoleToUser(req.UserID, req.RoleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
