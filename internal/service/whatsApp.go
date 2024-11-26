package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/linker/internal/mediator"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	SentMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error
	DeliveredMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error
	SetMessageStatus(ctx context.Context, message *entity.WhatsAppJSONReceived) error
	StartTemplate(ctx context.Context, template *dto.StartTemplateDto) error
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

func (w *whatsAppService) StartTemplate(ctx context.Context, template *dto.StartTemplateDto) error {
	hlog.Debug("whatsAppService.StartTemplate", fmt.Sprintf("%v", template))
	return w.whatsappMediator.StartTemplate(ctx, template)
}

func (w *whatsAppService) SentMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.SentMessage", fmt.Sprintf("%v", message))
	return w.whatsappMediator.SentMessage(ctx, message)
}

func (w *whatsAppService) DeliveredMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.DeliveredMessage", fmt.Sprintf("%v", message))
	return w.whatsappMediator.DeliveredMessage(ctx, message)
}

func (w *whatsAppService) SetMessageStatus(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.ReceiveMessage", fmt.Sprintf("%v", message))
	return w.whatsappMediator.SetStatus(ctx, message)
}
