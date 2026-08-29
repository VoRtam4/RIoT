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
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetUsers(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceUsers, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.LoadUsers()
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetUser(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceUsers, auth.OperationRead); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.LoadUserByUID(uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetMe(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal, ok := auth.PrincipalFromContext(c.Ctx)
	if !ok {
		sendError(c, msg.ID, "unauthorized")
		return
	}
	result := domainLogicLayer.LoadUserByID(principal.UserID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func UpdateUser(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceUsers, auth.OperationUpdate); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	input, err := parsePayload[graphQLModel.UserUpdateInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	result := domainLogicLayer.UpdateUser(uid, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func DisableUser(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceUsers, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	var reason *string
	if payload, ok := msg.Payload.(map[string]any); ok {
		if value, ok := payload["reason"].(string); ok {
			reason = &value
		}
	}
	result := domainLogicLayer.DisableUser(principal.UserID, uid, reason)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func EnableUser(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceUsers, auth.OperationUpdate); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.EnableUser(uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func RevokeUserSessions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceUsers, auth.OperationUpdate); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	if err := domainLogicLayer.RevokeUserSessions(uid); err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}

func GetSessions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceSessions, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.LoadUserSessions(principal.UserID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetSessionsByUser(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceUsers, auth.OperationRead); principal == nil {
		return
	}
	userUID, ok := payloadString(msg, "userUID")
	if !ok {
		sendError(c, msg.ID, "invalid userUID")
		return
	}
	result := domainLogicLayer.LoadUserSessionsByUserUID(userUID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func RevokeSession(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceSessions, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	if err := domainLogicLayer.RevokeSession(principal.UserID, uid, allowForeign); err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}

func RevokeOwnOtherSessions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceSessions, auth.OperationUpdate)
	if principal == nil {
		return
	}
	if principal.SessionID == nil {
		sendError(c, msg.ID, fmt.Errorf("current session not available").Error())
		return
	}
	if err := domainLogicLayer.RevokeOwnOtherSessions(principal.UserID, *principal.SessionID); err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}

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
