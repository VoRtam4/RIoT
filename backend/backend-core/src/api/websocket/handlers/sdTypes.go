/**
 * @file sdTypes.go
 * @brief WebSocket handlery pro práci s typy zdrojů dat.
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

func GetSDTypes(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceSDTypes, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetSDTypes()
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetSDType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceSDTypes, auth.OperationRead); principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	result := domainLogicLayer.GetSDType(uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func CreateSDType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceSDTypes, auth.OperationCreate); principal == nil {
		return
	}
	input, err := parsePayload[graphQLModel.SDTypeInput](msg)
	if err != nil {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	result := domainLogicLayer.CreateSDType(input)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func DeleteSDType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := AuthorizeOperation(c, msg, auth.ResourceSDTypes, auth.OperationDelete); principal == nil {
		return
	}
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid uid")
		return
	}
	if err := domainLogicLayer.DeleteSDType(uid); err != nil {
		sendError(c, msg.ID, err.Error())
		return
	}
	sendSuccess(c, msg.ID, nil)
}
