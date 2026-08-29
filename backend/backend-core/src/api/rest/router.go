/**
 * @file router.go
 * @brief Registrace REST endpointů Backend Core nad společným HTTP routerem.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
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
		r.Get("/sd-types/{uid}", handlers.GetSDType)
		r.Post("/sd-types", handlers.CreateSDType)
		r.Delete("/sd-types/{uid}", handlers.DeleteSDType)

		r.Get("/sd-instance", handlers.GetSDInstance)
		r.Get("/sd-instances", handlers.GetSDInstances)
		r.Get("/sd-instances/type/{uid}", handlers.GetSDInstancesByType)
		r.Get("/sd-instances/kpi/{uid}", handlers.GetSDInstancesByKpiDefinition)
		r.Patch("/sd-instances/{uid}", handlers.UpdateSDInstance)

		r.Get("/kpi-definitions", handlers.GetKPIDefinitions)
		r.Get("/kpi-definitions/type/{uid}", handlers.GetKPIDefinitionsBySDType)
		r.Get("/kpi-definitions/instance/{uid}", handlers.GetKPIDefinitionsBySDInstace)
		r.Get("/kpi-definitions/{uid}", handlers.GetKPIDefinition)
		r.Post("/kpi-definitions", handlers.CreateKPIDefinition)
		r.Put("/kpi-definitions/{uid}", handlers.UpdateKPIDefinition)
		r.Delete("/kpi-definitions/{uid}", handlers.DeleteKPIDefinition)

		r.Get("/raw/type/{uid}", handlers.GetRawDataPointsBySDType)
		r.Get("/raw", handlers.GetRawDataPoint)

		r.Get("/kpi-results", handlers.GetKPIResults)
		r.Get("/kpi-results/{uid}", handlers.GetKPIResultsByKPI)
		r.Post("/kpi-result", handlers.GetKPIResult)

		r.Get("/sd-instance-groups", handlers.GetSDInstanceGroups)
		r.Get("/sd-instance-groups/{uid}", handlers.GetSDInstanceGroup)
		r.Post("/sd-instance-groups", handlers.CreateSDInstanceGroup)
		r.Put("/sd-instance-groups/{uid}", handlers.UpdateSDInstanceGroup)
		r.Delete("/sd-instance-groups/{uid}", handlers.DeleteSDInstanceGroup)

		r.Get("/user-config", handlers.GetUserConfig)
		r.Post("/user-config", handlers.UpdateUserConfig)
		r.Delete("/user-config", handlers.DeleteUserConfig)

		r.Get("/users", handlers.GetUsers)
		r.Get("/users/me", handlers.GetMe)
		r.Get("/users/{uid}", handlers.GetUser)
		r.Patch("/users/{uid}", handlers.UpdateUser)
		r.Post("/users/{uid}/disable", handlers.DisableUser)
		r.Post("/users/{uid}/enable", handlers.EnableUser)
		r.Post("/users/{uid}/sessions/revoke", handlers.RevokeUserSessions)
		r.Get("/users/{uid}/sessions", handlers.GetSessionsByUser)
		r.Get("/users/{uid}/role", handlers.GetUserRole)
		r.Put("/users/{uid}/role", handlers.AssignRoleToUser)

		r.Get("/sessions", handlers.GetSessions)
		r.Delete("/sessions/{uid}", handlers.RevokeSession)
		r.Post("/sessions/revoke-others", handlers.RevokeOwnOtherSessions)

		r.Get("/roles", handlers.GetRoles)
		r.Get("/roles/{uid}", handlers.GetRoleByUID)
		r.Post("/roles", handlers.CreateRole)
		r.Put("/roles/{uid}", handlers.UpdateRole)
		r.Delete("/roles/{uid}", handlers.DeleteRole)
		r.Post("/roles/{uid}/clone", handlers.CloneRole)
		r.Get("/permissions", handlers.GetPermissions)
		r.Put("/permissions/{uid}/label", handlers.UpdatePermissionLabel)

		r.Get("/user-roles", handlers.GetRoles)
		r.Get("/user-roles/user", handlers.GetRole)
		r.Get("/user-roles/user/{uid}", handlers.GetUserRole)
		r.Put("/user-roles/user", handlers.AssignRoleToUser)

		r.Get("/api-keys", handlers.GetAPIKeys)
		r.Get("/api-keys/user/{uid}", handlers.GetAPIKeysByUser)
		r.Get("/api-keys/{uid}", handlers.GetAPIKey)
		r.Post("/api-keys", handlers.CreateAPIKey)
		r.Put("/api-keys/{uid}", handlers.UpdateAPIKey)
		r.Post("/api-keys/{uid}/revoke", handlers.RevokeAPIKey)
		r.Post("/api-keys/{uid}/rotate", handlers.RotateAPIKey)
		r.Put("/api-keys/{uid}/permissions", handlers.UpdateAPIKeyPermissions)
		r.Put("/api-keys/{uid}/restrictions", handlers.UpdateAPIKeyRestrictions)
		r.Delete("/api-keys/{uid}", handlers.DeleteAPIKey)

		r.Post("/time-series", handlers.ReadTimeSeries)
		r.Post("/time-series/distinct-tag-values", handlers.DistinctTimeSeriesTagValues)
		r.Post("/time-series/export", handlers.StartTimeSeriesExport)
		r.Get("/time-series/export/{uid}/status", handlers.GetTimeSeriesExport)
		r.Delete("/time-series/export/{uid}", handlers.CancelTimeSeriesExport)
		r.Post("/time-series/aggregate-kpi", handlers.ReadTimeSeriesAggregateKpi)
		r.Post("/time-series/export/aggregate-kpi", handlers.StartTimeSeriesExportAggregateKpi)
		r.Get("/time-series/export/{uid}", handlers.TimeSeriesExport)

		r.Get("/subscribe/{eventType}", handlers.SubscribeHandler)
	})
}
