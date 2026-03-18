package handlers

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
)

func GetKPIResults(c *websocket.Client, msg websocket.WebSocketMessage) {

	result := domainLogicLayer.GetKPIFulfillmentCheckResults()

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}
