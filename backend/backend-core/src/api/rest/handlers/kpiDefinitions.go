/**
 * @file kpiDefinitions.go
 * @brief REST handlery pro práci s KPI definicemi.
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

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/go-chi/chi/v5"
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
	uid := chi.URLParam(r, "uid")
	result := domainLogicLayer.GetKPIDefinition(principal.UserID, uid)
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
	uid := chi.URLParam(r, "uid")
	result := domainLogicLayer.GetKPIDefinitionsBySDType(principal.UserID, uid)
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
	uid := chi.URLParam(r, "uid")
	result := domainLogicLayer.GetKPIDefinitionsBySDInstance(principal.UserID, uid)
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
	uid := chi.URLParam(r, "uid")
	var input graphQLModel.KPIDefinitionInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	result := domainLogicLayer.UpdateKPIDefinition(principal.UserID, uid, input)
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
	uid := chi.URLParam(r, "uid")
	err := domainLogicLayer.DeleteKPIDefinition(principal.UserID, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
