package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetSDTypes(c *websocket.Client, msg websocket.WebSocketMessage) {

	result := domainLogicLayer.GetSDTypes()

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}

func GetSDType(c *websocket.Client, msg websocket.WebSocketMessage) {

	idFloat, ok := msg.Payload.(map[string]any)["id"].(float64)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: "invalid id"})
		return
	}

	result := domainLogicLayer.GetSDType(uint32(idFloat))

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}

func CreateSDType(c *websocket.Client, msg websocket.WebSocketMessage) {

	bytes, _ := json.Marshal(msg.Payload)

	var input graphQLModel.SDTypeInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: "invalid payload"})
		return
	}

	result := domainLogicLayer.CreateSDType(input)

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}

func DeleteSDType(c *websocket.Client, msg websocket.WebSocketMessage) {

	idFloat, ok := msg.Payload.(map[string]any)["id"].(float64)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: "invalid id"})
		return
	}

	err := domainLogicLayer.DeleteSDType(uint32(idFloat))
	if err != nil {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: err.Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true})
}
