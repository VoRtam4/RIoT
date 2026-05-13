/**
 * @file apiModel.go
 * @brief Sdílené modely pro obecnou API komunikaci přes WebSocket a návazné transportní vrstvy.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_commons
 */
package sharedModel

type MessageType string

const (
	MessageRequest     MessageType = "request"
	MessageResponse    MessageType = "response"
	MessageEvent       MessageType = "event"
	MessageSubscribe   MessageType = "subscribe"
	MessageUnsubscribe MessageType = "unsubscribe"
	MessageError       MessageType = "error"
)

type WebSocketMessage struct {
	Type    MessageType `json:"type"`
	ID      string      `json:"id,omitempty"`
	Action  string      `json:"action,omitempty"`
	Topic   string      `json:"topic,omitempty"`
	Payload any         `json:"payload,omitempty"`
	Success bool        `json:"success,omitempty"`
	Error   string      `json:"error,omitempty"`
}
