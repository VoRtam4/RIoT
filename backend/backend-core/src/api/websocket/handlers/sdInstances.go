package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetSDInstances(c *websocket.Client, msg websocket.WebSocketMessage) {

	result := domainLogicLayer.GetSDInstances()

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}

func UpdateSDInstance(c *websocket.Client, msg websocket.WebSocketMessage) {

	payload := msg.Payload.(map[string]any)

	idFloat, ok := payload["id"].(float64)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: "invalid id"})
		return
	}

	bytes, _ := json.Marshal(payload["input"])

	var input graphQLModel.SDInstanceUpdateInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: "invalid payload"})
		return
	}

	result := domainLogicLayer.UpdateSDInstance(uint32(idFloat), input)

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}
