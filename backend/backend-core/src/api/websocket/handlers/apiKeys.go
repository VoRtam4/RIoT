package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func GetAPIKeys(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := authorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationRead)
	if principal == nil {
		return
	}
	result := domainLogicLayer.LoadAPIKeysForUser(principal.UserID)
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

func CreateAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := authorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationCreate)
	if principal == nil {
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
	var input graphQLModel.CreateAPIKeyInput
	if err := json.Unmarshal(bytes, &input); err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	roleID, err := strconv.ParseUint(input.RoleID, 10, 32)
	if err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid roleId",
		})
		return
	}
	result := domainLogicLayer.CreateAPIKey(
		principal.UserID,
		uint32(roleID),
		input.Label,
		input.ExpiresAt,
	)
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
		Payload: map[string]string{"key": result.GetPayload()},
	})
}

func UpdateAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := authorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationUpdate)
	if principal == nil {
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
	idValue, ok := payload["id"]
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "missing id",
		})
		return
	}
	idFloat, ok := idValue.(float64)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	id := uint32(idFloat)
	apiKeyResult := domainLogicLayer.LoadAPIKeyByID(id)
	if apiKeyResult.IsFailure() {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: apiKeyResult.GetError().Error(),
		})
		return
	}
	if apiKeyResult.GetPayload().IsEmpty() {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "api key not found",
		})
		return
	}
	apiKey := apiKeyResult.GetPayload().GetPayload()
	if apiKey.UserID != principal.UserID {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	var input graphQLModel.UpdateAPIKeyInput
	if err := json.Unmarshal(bodyBytes, &input); err != nil {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid payload",
		})
		return
	}
	if input.Label != nil {
		apiKey.Label = *input.Label
	}
	if input.RoleID != nil {
		roleID, err := strconv.ParseUint(*input.RoleID, 10, 32)
		if err != nil {
			c.SafeSend(sharedModel.WebSocketMessage{
				Type:  sharedModel.MessageResponse,
				ID:    msg.ID,
				Error: "invalid roleId",
			})
			return
		}
		apiKey.RoleID = uint32(roleID)
	}
	if input.ExpiresAt != nil {
		apiKey.ExpiresAt = input.ExpiresAt
	}
	if input.Revoked != nil {
		apiKey.Revoked = *input.Revoked
	}
	if input.RateLimit != nil {
		apiKey.RateLimit = input.RateLimit
	}
	if err := domainLogicLayer.UpdateAPIKey(apiKey); err != nil {
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

func DeleteAPIKey(c *connection.Client, msg sharedModel.WebSocketMessage) {
	principal := authorizeOperation(c, msg, auth.ResourceAPIKeys, auth.OperationDelete)
	if principal == nil {
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
	idValue, ok := payload["id"]
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "missing id",
		})
		return
	}
	idFloat, ok := idValue.(float64)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "invalid id",
		})
		return
	}
	id := uint32(idFloat)
	apiKeyResult := domainLogicLayer.LoadAPIKeyByID(id)
	if apiKeyResult.IsFailure() {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: apiKeyResult.GetError().Error(),
		})
		return
	}
	if apiKeyResult.GetPayload().IsEmpty() {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "api key not found",
		})
		return
	}
	apiKey := apiKeyResult.GetPayload().GetPayload()
	if apiKey.UserID != principal.UserID {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return
	}
	if err := domainLogicLayer.DeleteAPIKey(id); err != nil {
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
