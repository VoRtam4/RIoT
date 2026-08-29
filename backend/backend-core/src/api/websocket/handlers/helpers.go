/**
 * @file helpers.go
 * @brief Pomocné funkce WebSocket handlerů pro autorizaci, parsování payloadů a odpovědi.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
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

func payloadString(msg sharedModel.WebSocketMessage, key string) (string, bool) {
	payload, ok := msg.Payload.(map[string]any)
	if !ok {
		return "", false
	}
	value, ok := payload[key].(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
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
