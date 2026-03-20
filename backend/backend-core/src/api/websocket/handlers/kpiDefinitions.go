package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetKPIDefinitions(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationRead); principal == nil {
		return
	}
	result := domainLogicLayer.GetKPIDefinitions()
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

func CreateKPIDefinition(c *connection.Client, msg sharedModel.WebSocketMessage) {
	if principal := authorizeOperation(c, msg, auth.ResourceKPIDefinitions, auth.OperationCreate); principal == nil {
		return
	}
	bytes, _ := json.Marshal(msg.Payload)
	var input graphQLModel.KPIDefinitionInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	result := domainLogicLayer.CreateKPIDefinition(input)
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
