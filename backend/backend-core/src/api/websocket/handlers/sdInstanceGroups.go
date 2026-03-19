package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func GetSDInstanceGroups(c *websocket.Client, msg websocket.WebSocketMessage) {
	if !auth.CanAccessOperation(c.Ctx, auth.ResourceSDInstances, auth.OperationRead) {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	result := domainLogicLayer.GetSDInstanceGroups()
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

func GetSDInstanceGroup(c *websocket.Client, msg websocket.WebSocketMessage) {
	if !auth.CanAccessOperation(c.Ctx, auth.ResourceSDInstances, auth.OperationRead) {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	idFloat, ok := payload["id"].(float64)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	result := domainLogicLayer.GetSDInstanceGroup(uint32(idFloat))
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

func CreateSDInstanceGroup(c *websocket.Client, msg websocket.WebSocketMessage) {
	if !auth.CanAccessOperation(c.Ctx, auth.ResourceSDTypes, auth.OperationCreate) {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	bytes, err := json.Marshal(msg.Payload)
	if err != nil {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	var input graphQLModel.SDInstanceGroupInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	result := domainLogicLayer.CreateSDInstanceGroup(input)
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

func UpdateSDInstanceGroup(c *websocket.Client, msg websocket.WebSocketMessage) {
	if !auth.CanAccessOperation(c.Ctx, auth.ResourceSDInstances, auth.OperationUpdate) {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	idFloat, ok := payload["id"].(float64)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	inputRaw, ok := payload["input"]
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "missing input",
		})
		return
	}
	bytes, err := json.Marshal(inputRaw)
	if err != nil {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid input",
		})
		return
	}
	var input graphQLModel.SDInstanceGroupInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid input",
		})
		return
	}
	result := domainLogicLayer.UpdateSDInstanceGroup(uint32(idFloat), input)
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

func DeleteSDInstanceGroup(c *websocket.Client, msg websocket.WebSocketMessage) {
	if !auth.CanAccessOperation(c.Ctx, auth.ResourceSDTypes, auth.OperationDelete) {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	idFloat, ok := payload["id"].(float64)
	if !ok {
		c.SafeSend(websocket.WebSocketMessage{
			Type:  websocket.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	err := domainLogicLayer.DeleteSDInstanceGroup(uint32(idFloat))
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
