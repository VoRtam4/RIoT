package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetUserConfig(c *websocket.Client, msg websocket.WebSocketMessage) {
	if !auth.CanAccessOperation(c.Ctx, auth.ResourceUserConfig, auth.OperationRead) {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	userID := c.Principal.UserID
	result := domainLogicLayer.GetUserConfig(userID)
	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: result.GetError().Error(),
		})
		return
	}
	c.SafeSend(websocket.WebSocketMessage{
		Type:    websocket.MessageResponse,
		ID:      msg.ID,
		Success: true,
		Payload: result.GetPayload(),
	})
}

func UpdateUserConfig(c *websocket.Client, msg websocket.WebSocketMessage) {
	if !auth.CanAccessOperation(c.Ctx, auth.ResourceUserConfig, auth.OperationUpdate) {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	userID := c.Principal.UserID
	bytes, _ := json.Marshal(msg.Payload)
	var input graphQLModel.UserConfigInput
	json.Unmarshal(bytes, &input)
	result := domainLogicLayer.UpdateUserConfig(userID, input)
	if result.IsFailure() {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: result.GetError().Error(),
		})
		return
	}
	c.SafeSend(websocket.WebSocketMessage{
		Type:    websocket.MessageResponse,
		ID:      msg.ID,
		Success: true,
		Payload: result.GetPayload(),
	})
}

func DeleteUserConfig(c *websocket.Client, msg websocket.WebSocketMessage) {
	if !auth.CanAccessOperation(c.Ctx, auth.ResourceUserConfig, auth.OperationUpdate) {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	userID := c.Principal.UserID
	err := domainLogicLayer.DeleteUserConfig(userID)
	if err != nil {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: err.Error(),
		})
		return
	}
	c.SafeSend(websocket.WebSocketMessage{
		Type:    websocket.MessageResponse,
		ID:      msg.ID,
		Success: true,
	})
}
