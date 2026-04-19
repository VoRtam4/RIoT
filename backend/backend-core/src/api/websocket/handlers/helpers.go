package handlers

import (
	"encoding/json"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func AuthorizeOperation(c *connection.Client, msg sharedModel.WebSocketMessage, operation string, opType string) *auth.Principal {
	principal, ok := auth.PrincipalFromContext(c.Ctx)
	if !ok {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "unauthorized",
		})
		return nil
	}
	if !auth.CanAccessOperation(principal, operation, opType) {
		c.SafeSend(sharedModel.WebSocketMessage{
			Type:  sharedModel.MessageResponse,
			ID:    msg.ID,
			Error: "forbidden",
		})
		return nil
	}
	return principal
}

func sendError(c *connection.Client, msgID string, err string) {
	c.SafeSend(sharedModel.WebSocketMessage{
		Type:  sharedModel.MessageResponse,
		ID:    msgID,
		Error: err,
	})
}

func sendSuccess(c *connection.Client, msgID string, payload any) {
	c.SafeSend(sharedModel.WebSocketMessage{
		Type:    sharedModel.MessageResponse,
		ID:      msgID,
		Success: true,
		Payload: payload,
	})
}

func parsePayload[T any](msg sharedModel.WebSocketMessage) (T, error) {
	var result T
	bytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(bytes, &result)
	return result, err
}

func parseID(payload map[string]any) (uint32, bool) {
	idFloat, ok := payload["id"].(float64)
	if !ok {
		return 0, false
	}
	return uint32(idFloat), true
}

func streamToClient[T any](c *connection.Client, topic string, subID string, ch <-chan T) {
	go func() {
		for payload := range ch {
			c.SafeSend(sharedModel.WebSocketMessage{
				Type:    sharedModel.MessageEvent,
				Topic:   topic,
				ID:      subID,
				Payload: payload,
			})
		}
	}()
}
