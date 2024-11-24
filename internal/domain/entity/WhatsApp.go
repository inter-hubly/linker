package entity

import (
	"time"

	"github.com/inter-hubly/pilot/domain/dto"
)

type whatsAppMessage struct {
	Body string `json:"body"`
}

func NewWhatsAppMessage(body string) *whatsAppMessage {
	return &whatsAppMessage{
		Body: body,
	}

}

func (w *whatsAppMessage) GetBody() string {
	return w.Body
}

func NormalizeWhatsAppMessage(message *dto.WhatsAppJSONReceived) *Chat {
	return &Chat{
		Id:        message.Id,
		MessageId: message.Metadata.MessageId,
		Message:   NewWhatsAppMessage(message.Metadata.Body),
		Delivered: ChatMessageTime{
			CreatedInDatabase: time.Now(),
			ReceivedAt:        message.Metadata.Timestamp,
		},
	}
}
