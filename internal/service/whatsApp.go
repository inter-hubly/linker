package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/linker/internal/mediator"
	"github.com/inter-hubly/linker/internal/repository"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	SentMessage(ctx context.Context, message *dto.SentTextDto) error
	DeliveredMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error
	SetMessageStatus(ctx context.Context, message *entity.WhatsAppJSONReceived) error
	StartTemplate(ctx context.Context, template *dto.StartTemplateDto) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppService
)

type whatsAppService struct {
	whatsappMediator   mediator.WhatsApp
	whatsappRepository repository.WhatsApp
}

func NewWhatsApp() *whatsAppService {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppService{
			whatsappMediator:   mediator.NewWhatsApp(),
			whatsappRepository: repository.NewWhatsApp(),
		}
	})
	return whatsApp
}

func (w *whatsAppService) StartTemplate(ctx context.Context, template *dto.StartTemplateDto) error {
	hlog.Debug("whatsAppService.StartTemplate", fmt.Sprintf("%v", template))
	return w.whatsappMediator.StartTemplate(ctx, template)
}

func (w *whatsAppService) SentMessage(ctx context.Context, template *dto.SentTextDto) error {
	hlog.Debug("whatsAppService.SentMessage", fmt.Sprintf("%v", template))
	message := entity.WhatsAppJSONReceived{
		Owner: entity.WhatsAppPhoneIdDto{
			PhoneNumberID:      "515719138282305",
			DisplayPhoneNumber: template.SenderAndReceiver.From,
		},
		SenderPhone: template.SenderAndReceiver.To,
		Metadata: entity.WhatsAppMetadataDto{
			Body: template.Message,
		},
	}

	return w.whatsappMediator.SentMessage(ctx, &message)
}

func (w *whatsAppService) DeliveredMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.DeliveredMessage", fmt.Sprintf("%v", message))

	return w.whatsappMediator.DeliveredMessage(ctx, message)
}

func (w *whatsAppService) SetMessageStatus(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.ReceiveMessage", fmt.Sprintf("%v", message))
	return w.whatsappMediator.SetStatus(ctx, message)
}
