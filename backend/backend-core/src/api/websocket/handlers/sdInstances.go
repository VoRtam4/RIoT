package handlers

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetSDInstances(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationRead); principal == nil {
		return
	}

	result := domainLogicLayer.GetSDInstances()
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}

	sendSuccess(c, msg.ID, result.GetPayload())
}

func UpdateSDInstance(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationUpdate); principal == nil {
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
	input, err := parsePayload[graphQLModel.SDInstanceUpdateInput](inputMsg)
	if err != nil {
		sendError(c, msg.ID, "invalid input")
		return
	}

	result := domainLogicLayer.UpdateSDInstance(id, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}

	sendSuccess(c, msg.ID, result.GetPayload())
}
