package websocket

import (
	"net/http"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigins := sharedUtils.NewSetFromSlice(strings.Split(sharedUtils.GetEnvironmentVariableValue("ALLOWED_ORIGINS").GetPayloadOrDefault("http://localhost:8080"), ","))
		return allowedOrigins.Contains(r.Header.Get("Origin"))
	},
}

func ServeWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authResult := auth.AuthenticateRequest(r)
		if authResult.IsFailure() {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if err := auth.ApplyAuthenticationResult(w, authResult.GetPayload()); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		client := NewClient(conn)
		client.UserID = authResult.GetPayload().UserID
		hub.Register(client)
		go client.writePump()
		go client.readPump(hub)
	}
}
