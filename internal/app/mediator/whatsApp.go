package mediator

import (
	"context"
	"time"

	dto "github.com/inter-hubly/linker/internal/app/domain/dto/pulse"
	dtoWhats "github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	"github.com/inter-hubly/linker/internal/app/domain/entity"
	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	SendMessage(ctx context.Context, message *dtoWhats.WhatsAppJSONReceived) error
	StartTemplate(ctx context.Context, template *dtoWhats.StartTemplateDto) error
	ReceiveMessage(ctx context.Context, received *dtoWhats.WhatsAppJSONReceived) error
}

type whatsAppMediator struct {
	messageProduct     string
	whatsAppRepository repository.WhatsApp
	whatsAppGateway    gateway.WhatsApp
	pulseGateway       gateway.Pulse
}

func NewWhatsApp() WhatsApp {
	return &whatsAppMediator{
		messageProduct:     "whatsapp",
		whatsAppRepository: repository.NewWhatsApp(),
		whatsAppGateway:    gateway.NewWhatsApp(),
		pulseGateway:       gateway.NewPulse(),
	}

}

func (w *whatsAppMediator) StartTemplate(ctx context.Context, template *dtoWhats.StartTemplateDto) error {
	message := dtoWhats.GatewayWhatsAppMessageDto{
		MessagingProduct: w.messageProduct,
		To:               template.SenderAndReceiver.To,
		Type:             dtoWhats.TemplateMessageType,
		Template: &dtoWhats.TemplateDto{
			Name: template.Name,
			LanguageDto: dtoWhats.LanguageDto{
				Code: template.Language,
			},
		},
	}
	res, err := w.whatsAppGateway.SendMessage(ctx, template.SenderAndReceiver.OwnerId, &message)
	// need persist even with errors
	chatDb := entity.Chat{}
	chatDb.OwnerId = template.SenderAndReceiver.OwnerId
	chatDb.Type = entity.ChatTemplate
	chatDb.TemplateName = template.Name
	chatDb.ToPhone = template.SenderAndReceiver.To
	chatDb.IsOwner = true

	if err == nil {
		chatDb.MessageId = res.Messages[0].Id

		chatDb.Audit = append(chatDb.Audit, entity.ChatMessageStatusTime{
			Status:     dtoWhats.StartStatus,
			ReceivedAt: time.Now().Unix(),
		})
	} else {
		hlog.Error("whatsAppMediator.StartTemplate", "error when send message to whatsApp", err)
		chatDb.Audit = append(chatDb.Audit, entity.ChatMessageStatusTime{
			Status:     dtoWhats.ErrorStatus,
			ReceivedAt: time.Now().Unix(),
		})
	}

	w.whatsAppRepository.PersistMessage(ctx, &chatDb)

	return nil
}

func (w *whatsAppMediator) SendMessage(ctx context.Context, message *dtoWhats.WhatsAppJSONReceived) error {

	messageToWhats := w.createTextMessage(ctx, message.Sender.PhoneNumber, message.Metadata.Body)

	whatsId, err := w.sendMessageToWhatsApp(ctx, message.Owner.PhoneNumberId, messageToWhats)

	if err != nil {
		return err
	}

	go func(id string) {
		w.persistMessageInElastic(ctx, id, message)
	}(whatsId)

	return nil
}

func (w *whatsAppMediator) persistMessageInElastic(ctx context.Context, whatsId string, received *dtoWhats.WhatsAppJSONReceived) error {
	chat := entity.Chat{
		MessageId: whatsId,
		Type:      entity.ChatText,
		OwnerId:   received.Owner.PhoneNumberId,
		ToPhone:   received.Sender.PhoneNumber,
		Message:   received.Metadata.Body,
		IsOwner:   true,
		Audit: []entity.ChatMessageStatusTime{
			{
				Status:     dtoWhats.StartStatus,
				ReceivedAt: time.Now().Unix(),
			},
		},
	}

	_, err := w.whatsAppRepository.PersistMessage(ctx, &chat)
	if err != nil {
		hlog.Error("whatsAppMediator.persistMessageInElastic", "error when persist message", err)
		return err
	}

	return nil
}

func (w *whatsAppMediator) ReceiveMessage(ctx context.Context, received *dtoWhats.WhatsAppJSONReceived) error {

	chat := entity.Chat{
		Type:        entity.ChatText,
		MessageId:   received.Metadata.MessageId,
		OwnerId:     received.Owner.PhoneNumberId,
		ToPhone:     received.Sender.PhoneNumber,
		Message:     received.Metadata.Body,
		ProfileName: received.Sender.ProfileName,
		IsOwner:     false,
	}
	_, err := w.whatsAppRepository.PersistMessage(ctx, &chat)

	if err != nil {
		return err
	}

	go func(entityChat *entity.Chat) {
		if err = w.pulseGateway.HandleMessage(ctx, received.Owner.PhoneNumberId, &dto.PulseDto{
			Message: entityChat.Message,
			ToPhone: entityChat.ToPhone,
		}); err != nil {
			hlog.Error("whatsAppMediator.ReceiveMessage", "error when persist message", err)
		}
	}(&chat)

	return nil
}

func (w *whatsAppMediator) sendMessageToWhatsApp(ctx context.Context, ownerId string, message *dtoWhats.GatewayWhatsAppMessageDto) (string, error) {
	resp, err := w.whatsAppGateway.SendMessage(ctx, ownerId, message)

	if err != nil {
		hlog.Error("whatsAppMediator.sendMessageToWhatsApp", "error when send message", err)
		return "", err
	}

	var whatsId string
	if resp != nil && resp.Messages != nil {
		whatsId = resp.Messages[0].Id
	}

	return whatsId, nil
}

const (
	SendError    = "SendError"
	PersistError = "PersistError"
)

type errValue struct {
	errType string
	err     error
}

func (w *whatsAppMediator) createTextMessage(ctx context.Context, to, body string) *dtoWhats.GatewayWhatsAppMessageDto {
	return &dtoWhats.GatewayWhatsAppMessageDto{
		MessagingProduct: w.messageProduct,
		To:               to,
		Type:             dtoWhats.TextMessageType,
		RecipientType:    "individual",
		Text: &dtoWhats.WhatsAppTextDto{
			PreviewUrl: false,
			Body:       body,
		},
	}
}

func (w *whatsAppMediator) createTemplateMessage(ctx context.Context, to, name string) *dtoWhats.GatewayWhatsAppMessageDto {
	// TODO Não seria melhor um DDD do numero de telefone?
	return &dtoWhats.GatewayWhatsAppMessageDto{
		MessagingProduct: w.messageProduct,
		To:               to,
		Type:             dtoWhats.TemplateMessageType,
		Template: &dtoWhats.TemplateDto{
			Name: name,
			LanguageDto: dtoWhats.LanguageDto{
				Code: "pt-br",
			},
		},
	}
}
