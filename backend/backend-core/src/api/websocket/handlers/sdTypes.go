package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetSDTypes(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDTypes, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetSDTypes()
	if result.IsFailure() {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: result.GetError().Error(),
		})
		return
	}
	c.SafeSend(sharedModel.WebSocketMessage{
		Type:    sharedModel.MessageResponse,
		ID:      msg.ID,
		Success: true,
		Payload: result.GetPayload(),
	})
}

func GetSDType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDTypes, auth.OperationRead); principal == nil {
		return
	}
	idFloat, ok := msg.Payload.(map[string]any)["id"].(float64)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	result := domainLogicLayer.GetSDType(uint32(idFloat))
	if result.IsFailure() {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: result.GetError().Error(),
		})
		return
	}
	c.SafeSend(sharedModel.WebSocketMessage{
		Type:    sharedModel.MessageResponse,
		ID:      msg.ID,
		Success: true,
		Payload: result.GetPayload(),
	})
}

func CreateSDType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDTypes, auth.OperationCreate); principal == nil {
		return
	}
	bytes, _ := json.Marshal(msg.Payload)
	var input graphQLModel.SDTypeInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	result := domainLogicLayer.CreateSDType(input)
	if result.IsFailure() {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: result.GetError().Error(),
		})
		return
	}
	c.SafeSend(sharedModel.WebSocketMessage{
		Type:    sharedModel.MessageResponse,
		ID:      msg.ID,
		Success: true,
		Payload: result.GetPayload(),
	})
}

func DeleteSDType(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDTypes, auth.OperationDelete); principal == nil {
		return
	}
	idFloat, ok := msg.Payload.(map[string]any)["id"].(float64)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	err := domainLogicLayer.DeleteSDType(uint32(idFloat))
	if err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: err.Error(),
		})
		return
	}
	c.SafeSend(sharedModel.WebSocketMessage{
		Type:    sharedModel.MessageResponse,
		ID:      msg.ID,
		Success: true,
	})
}
