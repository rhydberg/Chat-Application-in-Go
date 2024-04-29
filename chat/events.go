package chat

import "encoding/json"

const (
	EventSendMessage = "send_message"
)

type Event struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *Client) error

type SendMessageEvent struct {
	Message string `json:"message"`
	From string `json:"from"`
}