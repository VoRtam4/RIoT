package websocket

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/handlers"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

var requestHandlers = map[string]RequestHandler{
	"get_sd_types":   handlers.GetSDTypes,
	"get_sd_type":    handlers.GetSDType,
	"create_sd_type": handlers.CreateSDType,
	"delete_sd_type": handlers.DeleteSDType,

	"get_sd_instance":          handlers.GetSDInstance,
	"get_sd_instances":         handlers.GetSDInstances,
	"get_sd_instances_by_type": handlers.GetSDInstancesByType,
	"get_sd_instances_by_kpi":  handlers.GetSDInstancesByKpiDefinition,
	"update_sd_instance":       handlers.UpdateSDInstance,

	"get_kpi_definition":                 handlers.GetKPIDefinition,
	"get_kpi_definitions":                handlers.GetKPIDefinitions,
	"get_kpi_definitions_by_type":        handlers.GetKPIDefinitionsBySDType,
	"get_kpi_definitions_by_sd_instance": handlers.GetKPIDefinitionsBySDInstance,
	"create_kpi_definition":              handlers.CreateKPIDefinition,
	"update_kpi_definition":              handlers.UpdateKPIDefinition,
	"delete_kpi_definition":              handlers.DeleteKPIDefinition,

	"get_raw_by_sdtype":   handlers.GetRawDataPointsBySDType,
	"get_raw_by_instance": handlers.GetRawDataPoint,

	"get_kpi_results":        handlers.GetKPIResults,
	"get_kpi_results_by_kpi": handlers.GetKPIResultsByKPI,
	"get_kpi_result":         handlers.GetKPIResult,

	"get_sd_instance_groups":   handlers.GetSDInstanceGroups,
	"get_sd_instance_group":    handlers.GetSDInstanceGroup,
	"create_sd_instance_group": handlers.CreateSDInstanceGroup,
	"update_sd_instance_group": handlers.UpdateSDInstanceGroup,
	"delete_sd_instance_group": handlers.DeleteSDInstanceGroup,

	"get_user_config":    handlers.GetUserConfig,
	"update_user_config": handlers.UpdateUserConfig,
	"delete_user_config": handlers.DeleteUserConfig,

	"get_api_keys":   handlers.GetAPIKeys,
	"get_api_key":    handlers.GetAPIKey,
	"create_api_key": handlers.CreateAPIKey,
	"update_api_key": handlers.UpdateAPIKey,
	"delete_api_key": handlers.DeleteAPIKey,

	"get_user_roles":   handlers.GetRoles,
	"get_user_role":    handlers.GetUserRole,
	"get_role":         handlers.GetRole,
	"update_user_role": handlers.AssignRoleToUser,

	"time-series":                      handlers.StreamTimeSeries,
	"time-series-export":               handlers.StartTimeSeriesExport,
	"time-series-aggregate-kpi":        handlers.StreamTimeSeriesAggregateKpi,
	"time-series-export-aggregate-kpi": handlers.StartTimeSeriesExportAggregateKpi,
}

type RequestHandler func(c *connection.Client, msg sharedModel.WebSocketMessage)

func RouteMessage() func(h *connection.Hub, c *connection.Client, msg sharedModel.WebSocketMessage) {
	return func(h *connection.Hub, c *connection.Client, msg sharedModel.WebSocketMessage) {

		switch msg.Type {

		case sharedModel.MessageSubscribe:

			if msg.ID == "" {
				c.SafeSend(sharedModel.WebSocketMessage{
					Type:  sharedModel.MessageResponse,
					ID:    msg.ID,
					Error: "missing subscription id",
				})
				return
			}

			if oldClose, ok := c.Subscriptions[msg.ID]; ok {
				oldClose()
				delete(c.Subscriptions, msg.ID)
			}

			if err := subscribeClient(c, msg); err != nil {
				c.SafeSend(sharedModel.WebSocketMessage{
					Type:  sharedModel.MessageResponse,
					ID:    msg.ID,
					Error: err.Error(),
				})
			}

		case sharedModel.MessageUnsubscribe:

			if handlers.AuthorizeOperation(c, msg, auth.ResourceKPIResults, auth.OperationSubscribe) == nil {
				return
			}

			if closeFn, ok := c.Subscriptions[msg.ID]; ok {
				closeFn()
				delete(c.Subscriptions, msg.ID)
			}

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
