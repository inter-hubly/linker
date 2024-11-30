package mediator

import (
	"context"
	"errors"
	"fmt"
	"time"

	dto "github.com/inter-hubly/linker/internal/domain/dto/whatsapp"
	"github.com/inter-hubly/linker/internal/domain/entity"
	"github.com/inter-hubly/linker/internal/gateway"
	"github.com/inter-hubly/linker/internal/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
	StartTemplate(ctx context.Context, template *dto.StartTemplateDto) error
}

type whatsAppMediator struct {
	messageProduct     string
	whatsAppRepository repository.WhatsApp
	whatsAppGateway    gateway.WhatsApp
}

func NewWhatsApp() WhatsApp {
	return &whatsAppMediator{
		messageProduct:     "whatsapp",
		whatsAppRepository: repository.NewWhatsApp(),
		whatsAppGateway:    gateway.NewWhatsApp(),
	}

}

func (w *whatsAppMediator) StartTemplate(ctx context.Context, template *dto.StartTemplateDto) error {
	message := dto.GatewayWhatsAppMessageDto{
		MessagingProduct: w.messageProduct,
		To:               template.SenderAndReceiver.To,
		Type:             dto.TemplateMessageType,
		Template: &dto.TemplateDto{
			Name: template.Name,
			LanguageDto: dto.LanguageDto{
				Code: template.Language,
			},
		},
	}
	res, err := w.whatsAppGateway.SendMessage(ctx, template.SenderAndReceiver.OwnerNumberId, &message)
	// need persist even with errors
	chatDb := entity.Chat{}
	chatDb.OwnerId = template.SenderAndReceiver.OwnerNumberId
	chatDb.Type = entity.ChatTemplate
	chatDb.TemplateName = template.Name
	chatDb.OwnerPhone = template.SenderAndReceiver.From
	chatDb.ToPhoneId = template.SenderAndReceiver.To

	if err == nil {
		chatDb.MessageId = res.Messages[0].Id

		chatDb.Audit = append(chatDb.Audit, entity.ChatMessageStatusTime{
			Status:     dto.StartStatus,
			ReceivedAt: time.Now().String(),
		})
	} else {
		hlog.Error("whatsAppMediator.StartTemplate", "error when send message to whatsApp", err)
		chatDb.Audit = append(chatDb.Audit, entity.ChatMessageStatusTime{
			Status:     dto.ErrorStatus,
			ReceivedAt: fmt.Sprint(time.Now().Unix()),
		})
	}

	w.whatsAppRepository.PersistMessage(ctx, &chatDb)

	return nil
}

func (w *whatsAppMediator) SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {
	chanError := make(chan *errValue)

	messageToWhats := w.createTextMessage(ctx, message.SenderPhoneId, message.Metadata.Body)

	w.sendMessageToWhatsApp(ctx, messageToWhats, chanError)

	go w.persistMessageInElastic(ctx, message, chanError)

	for i := 0; i < 2; i++ {
		err := <-chanError
		if err != nil {
			if err.errType == SendError {
				return errors.New("error when sending mensage to whatsApp Gateway")
			}
		}
	}
	return nil
}

func (w *whatsAppMediator) persistMessageInElastic(ctx context.Context, received *dto.WhatsAppJSONReceived, chanError chan *errValue) {
	chat := entity.Chat{
		MessageId:  received.Metadata.MessageId,
		OwnerId:    received.Owner.PhoneNumberID,
		OwnerPhone: received.Owner.DisplayPhoneNumber,
		ToPhoneId:  received.SenderPhoneId,
		Message:    received.Metadata.Body,
		Audit: []entity.ChatMessageStatusTime{
			{
				Status: received.Status,
			},
		},
	}

	_, err := w.whatsAppRepository.PersistMessage(ctx, &chat)
	if err != nil {
		hlog.Error("whatsAppMediator.persistMessageInElastic", "error when persist message", err)
		chanError <- &errValue{
			errType: PersistError,
			err:     err,
		}
	}
	chanError <- nil
}

func (w *whatsAppMediator) sendMessageToWhatsApp(ctx context.Context, message *dto.GatewayWhatsAppMessageDto, chanError chan *errValue) {
	_, err := w.whatsAppGateway.SendMessage(ctx, "515719138282305", message)
	if err != nil {
		hlog.Error("whatsAppMediator.sendMessageToWhatsApp", "error when send message", err)
		chanError <- &errValue{
			errType: SendError,
			err:     err,
		}
	}
	chanError <- nil
}

const (
	SendError    = "SendError"
	PersistError = "PersistError"
)

type errValue struct {
	errType string
	err     error
}

func (w *whatsAppMediator) createTextMessage(ctx context.Context, to, body string) *dto.GatewayWhatsAppMessageDto {
	return &dto.GatewayWhatsAppMessageDto{
		MessagingProduct: w.messageProduct,
		To:               to,
		Type:             dto.TextMessageType,
		RecipientType:    "individual",
		Text: &dto.WhatsAppTextDto{
			PreviewUrl: false,
			Body:       body,
		},
	}
}

func (w *whatsAppMediator) createTemplateMessage(ctx context.Context, to, name string) *dto.GatewayWhatsAppMessageDto {
	// TODO Não seria melhor um DDD do numero de telefone?
	return &dto.GatewayWhatsAppMessageDto{
		MessagingProduct: w.messageProduct,
		To:               to,
		Type:             dto.TemplateMessageType,
		Template: &dto.TemplateDto{
			Name: name,
			LanguageDto: dto.LanguageDto{
				Code: "pt-br",
			},
		},
	}
}
