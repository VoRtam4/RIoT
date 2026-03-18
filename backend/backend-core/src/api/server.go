package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/graphql"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/rest"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

func StartServer() {
	r := chi.NewRouter()

	allowedOrigins := strings.Split(sharedUtils.GetEnvironmentVariableValue("ALLOWED_ORIGINS").GetPayloadOrDefault("http://localhost:8080"), ",")

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
	}).Handler)

	r.Get("/auth/login", auth.LoginHandler)
	r.Get("/auth/logout", auth.LogoutHandler)
	r.Get("/auth/callback", auth.CallbackHandler)

	rest.SetupRouter(r)

	hub := websocket.NewHub()
	go hub.StartEventListener()
	r.Get("/ws", websocket.ServeWebSocket(hub))

	r.Handle("/graphql", graphql.GetHandler())

	log.Println("Server running on :9090")
	log.Fatal(http.ListenAndServe(":9090", r))
}
