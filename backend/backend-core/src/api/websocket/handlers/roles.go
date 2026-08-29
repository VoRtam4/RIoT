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

func GetRoleByUID(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.LoadRoleByUID(uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetPermissions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.LoadPermissions()
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
	uid, ok := payload["uid"].(string)
	if !ok || uid == "" {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.LoadUserRoleByUID(uid)
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

func CreateRole(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationCreate); principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.RoleInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	result := domainLogicLayer.CreateRole(input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func UpdateRole(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationUpdate); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	input, err := parsePayload[graphQLModel.RoleInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	result := domainLogicLayer.UpdateRole(uid, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func DeleteRole(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationDelete); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	if err := domainLogicLayer.DeleteRole(uid); err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}

func CloneRole(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationCreate); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	label, ok := payloadString(msg, "label")
	if !ok {
		sendError(c, msg.ID, "invalid label")
		return
	}
	result := domainLogicLayer.CloneRole(uid, label)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func UpdatePermissionLabel(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceRoles, auth.OperationUpdate); principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	label, ok := payloadString(msg, "label")
	if !ok {
		sendError(c, msg.ID, "invalid label")
		return
	}
	result := domainLogicLayer.UpdatePermissionLabel(uid, label)
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
	err = domainLogicLayer.AssignRoleToUser(req.UserUID, req.RoleUID)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}
