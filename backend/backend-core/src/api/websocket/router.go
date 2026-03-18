package websocket

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/websocket/handlers"
)

type RequestHandler func(c *Client, msg WebSocketMessage)

var requestHandlers = map[string]RequestHandler{
	"get_sd_types":             handlers.GetSDTypes,
	"get_sd_type":              handlers.GetSDType,
	"create_sd_type":           handlers.CreateSDType,
	"delete_sd_type":           handlers.DeleteSDType,
	"get_sd_instances":         handlers.GetSDInstances,
	"update_sd_instance":       handlers.UpdateSDInstance,
	"get_kpi_definitions":      handlers.GetKPIDefinitions,
	"get_kpi_definition":       handlers.GetKPIDefinition,
	"create_kpi_definition":    handlers.CreateKPIDefinition,
	"update_kpi_definition":    handlers.UpdateKPIDefinition,
	"delete_kpi_definition":    handlers.DeleteKPIDefinition,
	"get_kpi_results":          handlers.GetKPIResults,
	"get_sd_instance_groups":   handlers.GetSDInstanceGroups,
	"get_sd_instance_group":    handlers.GetSDInstanceGroup,
	"create_sd_instance_group": handlers.CreateSDInstanceGroup,
	"update_sd_instance_group": handlers.UpdateSDInstanceGroup,
	"delete_sd_instance_group": handlers.DeleteSDInstanceGroup,
	"get_user_config":          handlers.GetUserConfig,
	"update_user_config":       handlers.UpdateUserConfig,
	//"delete_user_config":       handlers.DeleteUserConfig,
}

func RouteMessage(h *Hub, c *Client, msg WebSocketMessage) {
	switch msg.Type {
	case MessageSubscribe:
		c.subscriptions[msg.Topic] = true

	case MessageUnsubscribe:
		delete(c.subscriptions, msg.Topic)

	case MessageRequest:
		handler, ok := requestHandlers[msg.Action]
		if !ok {
			c.SafeSend(WebSocketMessage{Type: MessageResponse, ID: msg.ID, Error: "unknown action"})
			return
		}
		handler(c, msg)
	}
}
