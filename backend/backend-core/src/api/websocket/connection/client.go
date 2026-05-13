/**
 * @file client.go
 * @brief Model WebSocket klienta a bezpečné odesílání zpráv.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package connection

import (
	"context"
	"log"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn          *websocket.Conn
	Send          chan sharedModel.WebSocketMessage
	Ctx           context.Context
	Subscriptions map[string]func()
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn:          conn,
		Send:          make(chan sharedModel.WebSocketMessage, 32),
		Subscriptions: map[string]func(){},
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Println("WS write error:", err)
			return
		}
	}
}

func (c *Client) ReadPump(h *Hub, routeMessage func(h *Hub, c *Client, msg sharedModel.WebSocketMessage)) {
	defer func() {
		h.Unregister(c)
		c.Conn.Close()
	}()
	for {
		var msg sharedModel.WebSocketMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WS read error:", err)
			return
		}
		routeMessage(h, c, msg)
	}
}

func (c *Client) SafeSend(msg sharedModel.WebSocketMessage) {
	select {
	case c.Send <- msg:
	default:
	}
}
