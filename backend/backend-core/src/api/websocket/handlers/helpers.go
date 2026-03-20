package handlers

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func authorizeOperation(c *connection.Client, msg sharedModel.WebSocketMessage, operation string, opType string) *auth.Principal {
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
