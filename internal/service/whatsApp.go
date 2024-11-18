package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/linker/internal/mediator"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
	ReceiveMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppService
)

type whatsAppService struct {
	whatsappMediator mediator.WhatsApp
}

func NewWhatsApp() *whatsAppService {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppService{
			whatsappMediator: mediator.NewWhatsApp(),
		}
	})
	return whatsApp
}

func (w *whatsAppService) SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.SendMessage", fmt.Sprintf("%v", message))
	return w.whatsappMediator.Persist(ctx, message)
}

func (w *whatsAppService) ReceiveMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.ReceiveMessage", fmt.Sprintf("%v", message))
	return w.whatsappMediator.Persist(ctx, message)
}
