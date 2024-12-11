package service

import (
	"context"
	"fmt"
	"sync"

	dto "github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	"github.com/inter-hubly/linker/internal/app/mediator"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	ChangeStatusMessage(ctx context.Context, message *dto.ChangeStatusDto) error
	SendMessage(ctx context.Context, template *dto.SendTextDto) error
	StartTemplate(ctx context.Context, template *dto.StartTemplateDto) error
	ReceiveMessage(ctx context.Context, dto *dto.WhatsAppJSONReceived) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppService
)

type whatsAppService struct {
	whatsappMediator   mediator.WhatsApp
	whatsappRepository repository.WhatsApp
	clientRepository   repository.Client
}

func NewWhatsApp() *whatsAppService {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppService{
			whatsappMediator:   mediator.NewWhatsApp(),
			whatsappRepository: repository.NewWhatsApp(),
			clientRepository:   repository.NewClient(),
		}
	})
	return whatsApp
}

func (w *whatsAppService) StartTemplate(ctx context.Context, template *dto.StartTemplateDto) error {
	hlog.Debug("whatsAppService.StartTemplate", fmt.Sprintf("%v", template))
	return w.whatsappMediator.StartTemplate(ctx, template)
}

func (w *whatsAppService) SendMessage(ctx context.Context, template *dto.SendTextDto) error {
	hlog.Debug("whatsAppService.SendMessage", fmt.Sprintf("%v", template))

	message := dto.WhatsAppJSONReceived{
		Owner: dto.WhatsAppPhoneIdDto{
			PhoneNumberId: template.SenderAndReceiver.OwnerId,
		},
		Sender: dto.WhatsAppPhoneIdDto{
			PhoneNumberId: template.SenderAndReceiver.To,
		},
		Metadata: dto.WhatsAppMetadataDto{
			Body: template.Message,
		},
	}

	return w.whatsappMediator.SendMessage(ctx, &message)
}

func (w *whatsAppService) ReceiveMessage(ctx context.Context, dto *dto.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.ReceiveMessage", fmt.Sprintf("%v", dto))
	id, err := w.clientRepository.GetClientByPhoneId(ctx, dto.Sender.PhoneNumberId)
	if err != nil {
		hlog.Error("whatsAppService.ReceiveMessage", fmt.Sprintf("error geting number %s", dto.Sender.PhoneNumberId))
		return err
	}
	dto.Sender.PhoneNumber = id
	if err := w.whatsappMediator.ReceiveMessage(ctx, dto); err != nil {
		return err
	}

	return nil
}

func (w *whatsAppService) ChangeStatusMessage(ctx context.Context, message *dto.ChangeStatusDto) error {
	hlog.Debug("whatsAppService.ChangeStatusMessage", fmt.Sprintf("%v", message))

	if err := w.whatsappRepository.SetStatusMessageById(ctx, message.MessageId, message.Status, message.ExpirationTimeStamp); err != nil {
		return err
	}
	return nil
}
