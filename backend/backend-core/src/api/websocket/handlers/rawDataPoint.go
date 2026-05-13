/**
 * @file rawDataPoint.go
 * @brief WebSocket handlery pro čtení raw datových bodů.
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
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetRawDataPointsBySDType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceRawData, auth.OperationRead)
	if principal == nil {
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	sdTypeID, ok := parseID(payload)
	if !ok {
		sendError(c, msg.ID, "invalid sdTypeID")
		return
	}
	result := domainLogicLayer.GetRawDataPointsBySDType(sdTypeID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}

func GetRawDataPoint(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := AuthorizeOperation(c, msg, auth.ResourceRawData, auth.OperationRead)
	if principal == nil {
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		sendError(c, msg.ID, "invalid payload")
		return
	}
	sdInstanceID, ok := parseID(payload)
	if !ok {
		sendError(c, msg.ID, "invalid sdInstanceID")
		return
	}
	result := domainLogicLayer.GetRawDataPoint(sdInstanceID)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}
