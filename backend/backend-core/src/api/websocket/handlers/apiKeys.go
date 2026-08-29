/**
 * @file apiKeys.go
 * @brief WebSocket handlery pro správu API klíčů.
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

func GetAPIKeys(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.LoadAPIKeysForUser(principal.UserID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationRead)
	if principal == nil {
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
	result := domainLogicLayer.LoadAPIKeyByUID(principal.UserID, uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetAPIKeysByUser(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceUsers, auth.OperationRead); principal == nil {
		return
	}
	if principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationRead); principal == nil {
		return
	}
	userUID, ok := payloadString(msg, "userUID")
	if !ok {
		sendError(c, msg.ID, "invalid userUID")
		return
	}
	result := domainLogicLayer.LoadAPIKeysByUserUID(userUID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func CreateAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationCreate)
	if principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.APIKeyInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	result := domainLogicLayer.CreateAPIKey(principal.UserID, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, map[string]string{
		"key": result.GetPayload(),
	})
}

func UpdateAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
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
	input, err := parsePayload[graphQLModel.APIKeyInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	err = domainLogicLayer.UpdateAPIKeyForUser(principal.UserID, uid, input)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}

func RevokeAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.RevokeAPIKey(principal.UserID, uid, allowForeign)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func RotateAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.RotateAPIKey(principal.UserID, uid, allowForeign)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, map[string]string{
		"key": result.GetPayload(),
	})
}

func UpdateAPIKeyPermissions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	req, err := parsePayload[struct {
		PermissionUIDs []string `json:"permissionUIDs"`
	}](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.UpdateAPIKeyPermissions(principal.UserID, uid, req.PermissionUIDs, allowForeign)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func UpdateAPIKeyRestrictions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
		return
	}
	uid, ok := payloadString(msg, "uid")
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	input, err := parsePayload[graphQLModel.APIKeyRestrictionsInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	allowForeign := auth.CanAccessOperation(principal, auth.ResourceUsers, auth.OperationUpdate)
	result := domainLogicLayer.UpdateAPIKeyRestrictions(principal.UserID, uid, input, allowForeign)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func DeleteAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationDelete)
	if principal == nil {
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
	err := domainLogicLayer.DeleteAPIKeyForUser(principal.UserID, uid)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}
