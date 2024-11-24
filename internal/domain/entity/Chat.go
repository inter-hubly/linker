package entity

import "time"

type ChatType string

const WhatsApp ChatType = "WhatsApp"

type Chat struct {
	Id        string          `json:"id"`
	MessageId string          `json:"messageId"`
	Type      ChatType        `json:"type"`
	Received  ChatMessageTime `json:"received,omitempty"`
	Read      ChatMessageTime `json:"read,omitempty"`
	Delivered ChatMessageTime `json:"delivered,omitempty"`
	Message   ReceivedMessage `json:"message,omitempty"`
}

type ReceivedMessage interface {
	GetBody() string
}

type ChatMessageTime struct {
	CreatedInDatabase time.Time `json:"CreatedInDatabase"`
	ReceivedAt        string    `json:"ReceivedAt"`
}
