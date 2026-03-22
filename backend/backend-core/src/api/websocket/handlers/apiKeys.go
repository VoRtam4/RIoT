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

	id, ok := parseID(payload)
	if !ok {
		sendError(c, msg.ID, "invalid id")
		return
	}

	result := domainLogicLayer.LoadAPIKeyByID(principal.UserID, id)
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

	id, ok := parseID(payload)
	if !ok {
		sendError(c, msg.ID, "invalid id")
		return
	}

	input, err := parsePayload[graphQLModel.APIKeyInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}

	err = domainLogicLayer.UpdateAPIKeyForUser(principal.UserID, id, input)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}

	sendSuccess(c, msg.ID, nil)
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

	id, ok := parseID(payload)
	if !ok {
		sendError(c, msg.ID, "invalid id")
		return
	}

	err := domainLogicLayer.DeleteAPIKeyForUser(principal.UserID, id)
	if err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}

	sendSuccess(c, msg.ID, nil)
}
