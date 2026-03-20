package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetSDInstanceGroups(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetSDInstanceGroups()
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

func GetSDInstanceGroup(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationRead); principal == nil {
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	idFloat, ok := payload["id"].(float64)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	result := domainLogicLayer.GetSDInstanceGroup(uint32(idFloat))
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

func CreateSDInstanceGroup(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationCreate); principal == nil {
		return
	}
	bytes, err := json.Marshal(msg.Payload)
	if err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	var input graphQLModel.SDInstanceGroupInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	result := domainLogicLayer.CreateSDInstanceGroup(input)
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

func UpdateSDInstanceGroup(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationUpdate); principal == nil {
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	idFloat, ok := payload["id"].(float64)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	inputRaw, ok := payload["input"]
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "missing input",
		})
		return
	}
	bytes, err := json.Marshal(inputRaw)
	if err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid input",
		})
		return
	}
	var input graphQLModel.SDInstanceGroupInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid input",
		})
		return
	}
	result := domainLogicLayer.UpdateSDInstanceGroup(uint32(idFloat), input)
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

func DeleteSDInstanceGroup(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationDelete); principal == nil {
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	idFloat, ok := payload["id"].(float64)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	err := domainLogicLayer.DeleteSDInstanceGroup(uint32(idFloat))
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
