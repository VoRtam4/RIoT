/**
 * @file roles.go
 * @brief WebSocket handlery pro práci s rolemi a přiřazením rolí uživatelům.
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

func GetRoles(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.LoadRoles()
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetUserRole(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	userID, ok := parseID(payload)
	if !ok {
		sendError(c, msg.ID, "invalid userId")
		return
	}
	result := domainLogicLayer.LoadUserRole(userID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetRole(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.LoadUserRole(principal.UserID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func AssignRoleToUser(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationUpdate); principal == nil {
		return
	}
	req, err := parsePayload[graphQLModel.AssignRoleInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	err = domainLogicLayer.AssignRoleToUser(req.UserID, req.RoleID)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}
