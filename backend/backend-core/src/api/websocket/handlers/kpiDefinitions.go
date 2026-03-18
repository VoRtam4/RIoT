package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetKPIDefinitions(c *websocket.Client, msg websocket.WebSocketMessage) {

	result := domainLogicLayer.GetKPIDefinitions()

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}

func CreateKPIDefinition(c *websocket.Client, msg websocket.WebSocketMessage) {

	bytes, _ := json.Marshal(msg.Payload)

	var input graphQLModel.KPIDefinitionInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: "invalid payload"})
		return
	}

	result := domainLogicLayer.CreateKPIDefinition(input)

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}
