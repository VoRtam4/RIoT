/**
 * @file sdInstances.go
 * @brief WebSocket handlery pro práci s instancemi zdrojů dat.
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

func GetSDInstance(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationRead); principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.GetSDInstance(uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

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

func GetSDInstancesByType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationRead); principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.GetSDInstancesByType(uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetSDInstancesByKpiDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationRead)
	if principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.GetSDInstancesByKpiDefinition(principal.UserID, uid)
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
	uid, ok := parseUIDPayload(payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
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
	result := domainLogicLayer.UpdateSDInstance(uid, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}
