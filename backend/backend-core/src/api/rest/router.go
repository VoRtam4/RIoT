package rest

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/rest/handlers"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(r chi.Router) {
	r.Route("/rest", func(r chi.Router) {
		r.Use(auth.JWTAuthenticationMiddleware)

		r.Get("/sd-types", handlers.GetSDTypes)
		r.Get("/sd-types/{id}", handlers.GetSDType)
		r.Post("/sd-types", handlers.CreateSDType)
		r.Delete("/sd-types/{id}", handlers.DeleteSDType)

		r.Get("/sd-instance", handlers.GetSDInstance)
		r.Get("/sd-instances", handlers.GetSDInstances)
		r.Get("/sd-instances/type/{id}", handlers.GetSDInstancesByType)
		r.Get("/sd-instances/kpi/{id}", handlers.GetSDInstancesByKpiDefinition)
		r.Patch("/sd-instances/{id}", handlers.UpdateSDInstance)

		r.Get("/kpi-definitions", handlers.GetKPIDefinitions)
		r.Get("/kpi-definitions/{id}", handlers.GetKPIDefinition)
		r.Get("/kpi-definitions/type/{id}", handlers.GetKPIDefinitionsBySDType)
		r.Get("/kpi-definitions/instance/{id}", handlers.GetKPIDefinitionsBySDInstace)
		r.Post("/kpi-definitions", handlers.CreateKPIDefinition)
		r.Put("/kpi-definitions/{id}", handlers.UpdateKPIDefinition)
		r.Delete("/kpi-definitions/{id}", handlers.DeleteKPIDefinition)

		r.Get("/raw/{id}", handlers.GetRawDataPointsBySDType)
		r.Get("/raw", handlers.GetRawDataPoint)

		r.Get("/kpi-results", handlers.GetKPIResults)
		r.Get("/kpi-results/{id}", handlers.GetKPIResultsByKPI)
		r.Post("/kpi-result", handlers.GetKPIResult)

		r.Get("/sd-instance-groups", handlers.GetSDInstanceGroups)
		r.Get("/sd-instance-groups/{id}", handlers.GetSDInstanceGroup)
		r.Post("/sd-instance-groups", handlers.CreateSDInstanceGroup)
		r.Put("/sd-instance-groups/{id}", handlers.UpdateSDInstanceGroup)
		r.Delete("/sd-instance-groups/{id}", handlers.DeleteSDInstanceGroup)

		r.Get("/user-config", handlers.GetUserConfig)
		r.Post("/user-config", handlers.UpdateUserConfig)
		r.Delete("/user-config", handlers.DeleteUserConfig)

		r.Get("/user-roles", handlers.GetRoles)
		r.Get("/user-roles/user", handlers.GetRole)
		r.Get("/user-roles/user/{id}", handlers.GetUserRole)
		r.Put("/user-roles/user", handlers.AssignRoleToUser)

		r.Get("/api-keys", handlers.GetAPIKeys)
		r.Get("/api-keys/{id}", handlers.GetAPIKey)
		r.Post("/api-keys", handlers.CreateAPIKey)
		r.Put("/api-keys/{id}", handlers.UpdateAPIKey)
		r.Delete("/api-keys/{id}", handlers.DeleteAPIKey)

		r.Post("/time-series", handlers.ReadTimeSeries)
		r.Post("/time-series/distinct-tag-values", handlers.DistinctTimeSeriesTagValues)
		r.Post("/time-series/export", handlers.StartTimeSeriesExport)
		r.Post("/time-series/aggregate-kpi", handlers.ReadTimeSeriesAggregateKpi)
		r.Post("/time-series/export/aggregate-kpi", handlers.StartTimeSeriesExportAggregateKpi)
		r.Get("/time-series/export/{id}", handlers.TimeSeriesExport)

		r.Get("/subscribe/{eventType}", handlers.SubscribeHandler)
	})
}
