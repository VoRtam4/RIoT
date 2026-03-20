package websocket

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/websocket/handlers"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

type RequestHandler func(c *connection.Client, msg sharedModel.WebSocketMessage)

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
	"delete_user_config":       handlers.DeleteUserConfig,
	"get_api_keys":             handlers.GetAPIKeys,
	"create_api_key":           handlers.CreateAPIKey,
	"update_api_key":           handlers.UpdateAPIKey,
	"delete_api_key":           handlers.DeleteAPIKey,
}

func RouteMessage() func(h *connection.Hub, c *connection.Client, msg sharedModel.WebSocketMessage) {
	return func(h *connection.Hub, c *connection.Client, msg sharedModel.WebSocketMessage) {
		switch msg.Type {
		case sharedModel.MessageSubscribe:
			if principle := handlers.authorizeOperation(c, msg, auth.ResourceEvents, auth.OperationSubscribe); principle == nil {
				return
			}
			c.Subscriptions[msg.Topic] = true

		case sharedModel.MessageUnsubscribe:
			if principle := handlers.authorizeOperation(c, msg, auth.ResourceEvents, auth.OperationSubscribe); principle == nil {
				return
			}
			delete(c.Subscriptions, msg.Topic)

		case sharedModel.MessageRequest:
			handler, ok := requestHandlers[msg.Action]
			if !ok {
				c.SafeSend(sharedModel.WebSocketMessage{
					Type:  sharedModel.MessageResponse,
					ID:    msg.ID,
					Error: "unknown action",
				})
				return
			}
			handler(c, msg)
		}
	}
}
