/**
 * @file user.go
 * @brief WebSocket handlery pro práci s uživatelskou konfigurací.
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
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetUserConfig(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceUserConfig, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.GetUserConfig(principal.UserID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func UpdateUserConfig(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceUserConfig, auth.OperationUpdate)
	if principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.UserConfigInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	result := domainLogicLayer.UpdateUserConfig(principal.UserID, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func DeleteUserConfig(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceUserConfig, auth.OperationDelete)
	if principal == nil {
		return
	}
	if err := domainLogicLayer.DeleteUserConfig(principal.UserID); err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}
