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

		r.Get("/sd-instances", handlers.GetSDInstances)
		r.Patch("/sd-instances/{id}", handlers.UpdateSDInstance)

		r.Get("/kpi-definitions", handlers.GetKPIDefinitions)
		r.Get("/kpi-definitions/{id}", handlers.GetKPIDefinition)
		r.Post("/kpi-definitions", handlers.CreateKPIDefinition)
		r.Put("/kpi-definitions/{id}", handlers.UpdateKPIDefinition)
		r.Delete("/kpi-definitions/{id}", handlers.DeleteKPIDefinition)

		r.Get("/kpi-results", handlers.GetKPIResults)

		r.Get("/sd-instance-groups", handlers.GetSDInstanceGroups)
		r.Get("/sd-instance-groups/{id}", handlers.GetSDInstanceGroup)
		r.Post("/sd-instance-groups", handlers.CreateSDInstanceGroup)
		r.Put("/sd-instance-groups/{id}", handlers.UpdateSDInstanceGroup)
		r.Delete("/sd-instance-groups/{id}", handlers.DeleteSDInstanceGroup)

		r.Get("/user-config", handlers.GetUserConfig)
		r.Post("/user-config", handlers.UpdateUserConfig)
		r.Delete("/user-config", handlers.DeleteUserConfig)

		r.Get("/api-keys", handlers.GetAPIKeys)
		r.Get("/api-keys/{id}", handlers.GetAPIKey)
		r.Post("/api-keys", handlers.CreateAPIKey)
		r.Put("/api-keys/{id}", handlers.UpdateAPIKey)
		r.Delete("/api-keys/{id}", handlers.DeleteAPIKey)

		r.Get("/subscribe", handlers.EventsStream)
	})
}
