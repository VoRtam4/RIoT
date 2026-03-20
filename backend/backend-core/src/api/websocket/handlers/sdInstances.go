package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetSDInstances(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetSDInstances()
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

func UpdateSDInstance(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceSDInstances, auth.OperationUpdate); principal == nil {
		return
	}
	payload := msg.Payload.(map[string]any)
	idFloat, ok := payload["id"].(float64)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	bytes, _ := json.Marshal(payload["input"])
	var input graphQLModel.SDInstanceUpdateInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	result := domainLogicLayer.UpdateSDInstance(uint32(idFloat), input)
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
