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
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid sdTypeUID")
		return
	}
	result := domainLogicLayer.GetRawDataPointsBySDTypeUID(uid)
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
	uid, ok := parseUIDPayload(msg.Payload)
	if !ok {
		sendError(c, msg.ID, "invalid sdInstanceUID")
		return
	}
	result := domainLogicLayer.GetRawDataPointBySDInstanceUID(uid)
	if result.IsFailure() {
		sendError(c, msg.ID, result.GetError().Error())
		return
	}
	sendSuccess(c, msg.ID, result.GetPayload())
}
