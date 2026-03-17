package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetUserConfig(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Context().Value(auth.UserIdContextIdentifier).(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	result := domainLogicLayer.GetUserConfig(uint32(userID))

	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.GetPayload())
}

func UpdateUserConfig(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Context().Value(auth.UserIdContextIdentifier).(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	var input graphQLModel.UserConfigInput

	json.NewDecoder(r.Body).Decode(&input)

	result := domainLogicLayer.UpdateUserConfig(uint32(userID), input)

	if result.IsFailure() {
		http.Error(w, result.GetError().Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result.GetPayload())
}

func DeleteUserConfig(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.Context().Value(auth.UserIdContextIdentifier).(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	err := domainLogicLayer.DeleteUserConfig(uint32(userID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
