package websocket

import (
	"sync"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
)

type Hub struct {
	clients map[*Client]bool
	mutex   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: map[*Client]bool{},
	}
}

func (h *Hub) Register(c *Client) {
	h.mutex.Lock()
	h.clients[c] = true
	h.mutex.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.mutex.Lock()
	delete(h.clients, c)
	close(c.send)
	h.mutex.Unlock()
}

func (h *Hub) StartEventListener() {
	sub := events.GetEventBus().Subscribe([]events.EventType{
		events.SDInstanceRegisteredEventType,
		events.KPIFulfillmentCheckedEventType,
	}, 128)
	go func() {
		for event := range sub.Channel {
			h.mutex.RLock()
			for client := range h.clients {
				if !client.subscriptions[string(event.Type)] {
					continue
				}
				client.SafeSend(WebSocketMessage{
					Type:    MessageEvent,
					Topic:   string(event.Type),
					Payload: event.Payload,
				})
			}
			h.mutex.RUnlock()
		}
	}()
}
