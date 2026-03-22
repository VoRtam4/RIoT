package handlers

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetKPIDefinitions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationRead); principal == nil {
		return
	}

	result := domainLogicLayer.GetKPIDefinitions()
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}

	sendSuccess(c, msg.ID, result.GetPayload())
}

func CreateKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationCreate); principal == nil {
		return
	}

	input, err := parsePayload[graphQLModel.KPIDefinitionInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}

	result := domainLogicLayer.CreateKPIDefinition(input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}

	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationRead); principal == nil {
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

	result := domainLogicLayer.GetKPIDefinition(id)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}

	sendSuccess(c, msg.ID, result.GetPayload())
}

func UpdateKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationUpdate); principal == nil {
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

	inputRaw, ok := payload["input"]
	if !ok {
		sendError(c, msg.ID, "missing input")
		return
	}

	inputMsg := sharedModel.WebSocketMessage{Payload: inputRaw}
	input, err := parsePayload[graphQLModel.KPIDefinitionInput](inputMsg)
	if err != nil {
		sendError(c, msg.ID, "invalid input")
		return
	}

	result := domainLogicLayer.UpdateKPIDefinition(id, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}

	sendSuccess(c, msg.ID, result.GetPayload())
}

func DeleteKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationDelete); principal == nil {
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

	if err := domainLogicLayer.DeleteKPIDefinition(id); err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}

	sendSuccess(c, msg.ID, nil)
}
