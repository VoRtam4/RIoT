package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan WebSocketMessage

	subscriptions map[string]bool

	UserID string
	roles  []string
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:          conn,
		send:          make(chan WebSocketMessage, 32),
		subscriptions: map[string]bool{},
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Println("WS write error:", err)
			return
		}
	}
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister(c)
		c.conn.Close()
	}()
	for {
		var msg WebSocketMessage
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WS read error:", err)
			return
		}
		RouteMessage(h, c, msg)
	}
}

func (c *Client) SafeSend(msg WebSocketMessage) {
	select {
	case c.send <- msg:
	default:
	}
}
