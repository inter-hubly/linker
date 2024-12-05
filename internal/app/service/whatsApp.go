package service

import (
	"context"
	"fmt"
	"sync"

	dto2 "github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	"github.com/inter-hubly/linker/internal/app/domain/entity"
	"github.com/inter-hubly/linker/internal/app/mediator"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	ChangeStatusMessage(ctx context.Context, message *dto2.ChangeStatusDto) error
	SendMessage(ctx context.Context, template *dto2.SendTextDto) error
	StartTemplate(ctx context.Context, template *dto2.StartTemplateDto) error
	ReceiveMessage(ctx context.Context, dto *dto2.WhatsAppJSONReceived) error
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

func (w *whatsAppService) StartTemplate(ctx context.Context, template *dto2.StartTemplateDto) error {
	hlog.Debug("whatsAppService.StartTemplate", fmt.Sprintf("%v", template))
	return w.whatsappMediator.StartTemplate(ctx, template)
}

func (w *whatsAppService) SendMessage(ctx context.Context, template *dto2.SendTextDto) error {
	hlog.Debug("whatsAppService.SendMessage", fmt.Sprintf("%v", template))
	message := dto2.WhatsAppJSONReceived{
		Owner: dto2.WhatsAppPhoneIdDto{
			PhoneNumberID:      template.SenderAndReceiver.OwnerNumberId,
			DisplayPhoneNumber: template.SenderAndReceiver.From,
		},
		SenderPhoneId: template.SenderAndReceiver.To,
		Metadata: dto2.WhatsAppMetadataDto{
			Body: template.Message,
		},
	}

	return w.whatsappMediator.SendMessage(ctx, &message)
}

func (w *whatsAppService) ReceiveMessage(ctx context.Context, dto *dto2.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppService.ReceiveMessage", fmt.Sprintf("%v", dto))
	chat := entity.Chat{
		MessageId:  dto.Metadata.MessageId,
		OwnerId:    dto.Owner.PhoneNumberID,
		OwnerPhone: dto.Owner.DisplayPhoneNumber,
		ToPhoneId:  dto.SenderPhoneId,
		Message:    dto.Metadata.Body,
	}
	_, err := w.whatsappRepository.PersistMessage(ctx, &chat)

	if err != nil {
		return err
	}
	return nil
}

func (w *whatsAppService) ChangeStatusMessage(ctx context.Context, message *dto2.ChangeStatusDto) error {
	hlog.Debug("whatsAppService.ChangeStatusMessage", fmt.Sprintf("%v", message))

	if err := w.whatsappRepository.SetStatusMessageById(ctx, message.MessageId, message.Status); err != nil {
		return err
	}
	return nil
}
