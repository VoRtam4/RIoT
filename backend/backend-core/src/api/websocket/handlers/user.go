package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetUserConfig(c *connection.Client, msg sharedModel.WebSocketMessage) {
	var principal *auth.Principal
	if principal = authorizeOperation(c, msg, auth.ResourceUserConfig, auth.OperationRead); principal == nil {
		return
	}
	userID := principal.UserID
	result := domainLogicLayer.GetUserConfig(userID)
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

func UpdateUserConfig(c *connection.Client, msg sharedModel.WebSocketMessage) {
	var principal *auth.Principal
	if principal = authorizeOperation(c, msg, auth.ResourceUserConfig, auth.OperationUpdate); principal == nil {
		return
	}
	userID := principal.UserID
	bytes, _ := json.Marshal(msg.Payload)
	var input graphQLModel.UserConfigInput
	json.Unmarshal(bytes, &input)
	result := domainLogicLayer.UpdateUserConfig(userID, input)
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

func DeleteUserConfig(c *connection.Client, msg sharedModel.WebSocketMessage) {
	var principal *auth.Principal
	if principal = authorizeOperation(c, msg, auth.ResourceUserConfig, auth.OperationDelete); principal == nil {
		return
	}
	err := domainLogicLayer.DeleteUserConfig(principal.UserID)
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
