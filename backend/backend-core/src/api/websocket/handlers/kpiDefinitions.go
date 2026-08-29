/**
 * @file kpiDefinitions.go
 * @brief WebSocket handlery pro práci s KPI definicemi.
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

func GetKPIDefinitions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.GetKPIDefinitions(principal.UserID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func CreateKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationCreate)
	if principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.KPIDefinitionInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	result := domainLogicLayer.CreateKPIDefinition(principal.UserID, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetKPIDefinitionsBySDType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationRead)
	if principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.GetKPIDefinitionsBySDType(principal.UserID, uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetKPIDefinitionsBySDInstance(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationRead)
	if principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.GetKPIDefinitionsBySDInstance(principal.UserID, uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationRead)
	if principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.GetKPIDefinition(principal.UserID, uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func UpdateKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationUpdate)
	if principal == nil {
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
	input, err := parsePayload[graphQLModel.KPIDefinitionInput](inputMsg)
	if err != nil {
		sendError(c, msg.ID, "invalid input")
		return
	}
	result := domainLogicLayer.UpdateKPIDefinition(principal.UserID, uid, input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func DeleteKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationDelete)
	if principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	if err := domainLogicLayer.DeleteKPIDefinition(principal.UserID, uid); err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}

func parseUIDPayload(payload any) (string, bool) {
	if uid, ok := payload.(string); ok {
		return uid, uid != ""
	}
	payloadMap, ok := payload.(map[string]any)
	if !ok {
		return "", false
	}
	uid, ok := payloadMap["uid"].(string)
	return uid, ok && uid != ""
}
