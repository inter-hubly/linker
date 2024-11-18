package gateway

import (
	"context"
	"sync"

	"github.com/inter-hubly/linker/internal/domain/dto"
)

type WhatsApp interface {
	SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
	ReceiveMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppGateway
)

type whatsAppGateway struct {
}

func NewWhatsApp() *whatsAppGateway {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppGateway{}
	})
	return whatsApp
}
