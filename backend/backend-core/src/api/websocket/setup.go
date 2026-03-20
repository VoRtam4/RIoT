package websocket

import (
	"net/http"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/websocket/connection"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
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

func EventListener() func(*connection.Hub) {
	return func(h *connection.Hub) {
		sub := events.GetEventBus().Subscribe([]events.EventType{
			events.SDInstanceRegisteredEventType,
			events.KPIFulfillmentCheckedEventType,
		}, 128)
		go func() {
			for event := range sub.Channel {
				h.Mutex.RLock()
				for client := range h.Clients {
					if !client.Subscriptions[string(event.Type)] {
						continue
					}
					client.SafeSend(sharedModel.WebSocketMessage{
						Type:    sharedModel.MessageEvent,
						Topic:   string(event.Type),
						Payload: event.Payload,
					})
				}
				h.Mutex.RUnlock()
			}
		}()
	}
}
