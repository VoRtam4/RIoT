package rest

import (
	"net/http"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/rest/handlers"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/go-chi/chi/v5"
)

func SetupRouter() http.Handler {

	r := chi.NewRouter()

	r.Use(auth.JWTAuthenticationMiddleware)

	r.Route("/api", func(r chi.Router) {

		r.Get("/sd-types", handlers.GetSDTypes)
		r.Get("/sd-types/{id}", handlers.GetSDType)
		r.Post("/sd-types", handlers.CreateSDType)
		r.Delete("/sd-types/{id}", handlers.DeleteSDType)

		r.Get("/sd-instance-groups", handlers.GetSDInstanceGroups)
		r.Get("/sd-instance-groups/{id}", handlers.GetSDInstanceGroup)
		r.Post("sd-instance-groups", handlers.CreateSDInstanceGroup)
		r.Put("/sd-instance-groups/{id}", handlers.UpdateSDInstanceGroup)
		r.Delete("/sd-instance-groups/{id}", handlers.DeleteSDInstanceGroup)

		r.Get("/user-config", handlers.GetUserConfig)
		r.Post("/user-config", handlers.UpdateUserConfig)
		r.Delete("/user-config", handlers.DeleteUserConfig)

		r.Get("/api/sd-instances", handlers.GetSDInstances)

		r.Get("/api/kpi-definitions", handlers.GetKPIDefinitions)

		r.Get("/kpi-results", handlers.GetKPIResults)
	})

	return r
}
