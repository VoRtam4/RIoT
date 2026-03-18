package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetUserConfig(c *websocket.Client, msg websocket.WebSocketMessage) {

	userID, _ := strconv.ParseUint(c.UserID, 10, 32)

	result := domainLogicLayer.GetUserConfig(uint32(userID))

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}

func UpdateUserConfig(c *websocket.Client, msg websocket.WebSocketMessage) {

	userID, _ := strconv.ParseUint(c.UserID, 10, 32)

	bytes, _ := json.Marshal(msg.Payload)

	var input graphQLModel.UserConfigInput
	json.Unmarshal(bytes, &input)

	result := domainLogicLayer.UpdateUserConfig(uint32(userID), input)

	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: result.GetError().Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true, Payload: result.GetPayload()})
}

func DeleteUserConfig(c *websocket.Client, msg websocket.WebSocketMessage) {

	userID, _ := strconv.ParseUint(c.UserID, 10, 32)

	err := domainLogicLayer.DeleteUserConfig(uint32(userID))

	if err == nil {
		c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Error: err.Error()})
		return
	}

	c.SafeSend(websocket.WebSocketMessage{Type: websocket.MessageResponse, ID: msg.ID, Success: true})
}
