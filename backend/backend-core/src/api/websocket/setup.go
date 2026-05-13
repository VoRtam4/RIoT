/**
 * @file setup.go
 * @brief Inicializace WebSocket endpointu a navázání připojení na autentizovaného principala.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package websocket

import (
	"net/http"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
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

func ServeWebSocket(hub *connection.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		principal, err := auth.AuthenticatePrincipal(w, r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := auth.ContextWithPrincipal(r.Context(), principal)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		client := connection.NewClient(conn)
		client.Ctx = ctx
		hub.Register(client)
		go client.WritePump()
		go client.ReadPump(hub, RouteMessage())
	}
}
