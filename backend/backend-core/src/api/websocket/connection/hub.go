package connection

import (
	"sync"
)

type Hub struct {
	Clients map[*Client]bool
	Mutex   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: map[*Client]bool{},
	}
}

func (h *Hub) Register(c *Client) {
	h.Mutex.Lock()
	h.Clients[c] = true
	h.Mutex.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.Mutex.Lock()
	delete(h.Clients, c)
	h.Mutex.Unlock()
	for _, closeFn := range c.Subscriptions {
		closeFn()
	}
	close(c.Send)
}

func (h *Hub) StartEventListener(listener func(h *Hub)) {
	listener(h)
}
