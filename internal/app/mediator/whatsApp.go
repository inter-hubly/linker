package mediator

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/linker/internal/app/domain/entity"
	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/server"
)

type WhatsApp interface {
	SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
	StartTemplate(ctx context.Context, campaignId uuid.UUID, template *dto.GatewayWhatsAppMessageDto) error
	ReceiveMessage(ctx context.Context, received *dto.WhatsAppJSONReceived) error
}

type whatsAppMediator struct {
	messageProduct     string
	whatsAppRepository repository.WhatsApp
	whatsAppGateway    gateway.WhatsApp
	pulseGateway       gateway.Pulse
}

func NewWhatsApp() WhatsApp {
	var whatsAppGateway gateway.WhatsApp
	if server.GetEnvironment().Env != "development" {
		whatsAppGateway = gateway.NewWhatsAppMock()
	} else {
		whatsAppGateway = gateway.NewWhatsApp()
	}

	return &whatsAppMediator{
		messageProduct:     "whatsapp",
		whatsAppRepository: repository.NewWhatsApp(),
		whatsAppGateway:    whatsAppGateway,
		pulseGateway:       gateway.NewPulse(),
	}

}

func (w *whatsAppMediator) StartTemplate(ctx context.Context, campaignId uuid.UUID, template *dto.GatewayWhatsAppMessageDto) error {
	tenantId := hctx.Tenant.Get(ctx)

	res, err := w.whatsAppGateway.SendMessage(ctx, tenantId, template)
	// need persist even with errors
	chatDb := entity.Chat{}
	chatDb.OwnerId = tenantId
	chatDb.Type = entity.ChatTemplate
	chatDb.CampaignId = campaignId
	chatDb.ToPhoneId = template.To
	chatDb.IsOwner = true

	if err == nil {
		chatDb.MessageId = res.Messages[0].Id

		chatDb.Audit = append(chatDb.Audit, entity.ChatMessageStatusTime{
			Status:     dto.StartStatus,
			ReceivedAt: time.Now().Unix(),
		})
	} else {
		hlog.Error(ctx, "whatsAppMediator.StartTemplate", fmt.Sprint("error when send message to whatsApp", err))
		chatDb.Audit = append(chatDb.Audit, entity.ChatMessageStatusTime{
			Status:     dto.ErrorStatus,
			ReceivedAt: time.Now().Unix(),
		})
	}

	w.whatsAppRepository.PersistMessage(ctx, &chatDb)

	return nil
}

func (w *whatsAppMediator) SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {

	messageToWhats := w.createTextMessage(ctx, message.Sender.PhoneNumberId, message.Metadata.Body)

	whatsId, err := w.sendMessageToWhatsApp(ctx, message.Owner.PhoneNumberId, messageToWhats)

	if err != nil {
		return err
	}

	go func(id string) {
		w.persistMessageInElastic(ctx, id, message)
	}(whatsId)

	return nil
}

func (w *whatsAppMediator) persistMessageInElastic(ctx context.Context, whatsId string, received *dto.WhatsAppJSONReceived) error {
	chat := entity.Chat{
		MessageId: whatsId,
		Type:      entity.ChatText,
		OwnerId:   received.Owner.PhoneNumberId,
		ToPhoneId: received.Sender.PhoneNumberId,
		Message:   received.Metadata.Body,
		IsOwner:   true,
		Audit: []entity.ChatMessageStatusTime{
			{
				Status:     dto.StartStatus,
				ReceivedAt: time.Now().Unix(),
			},
		},
	}

	_, err := w.whatsAppRepository.PersistMessage(ctx, &chat)
	if err != nil {
		hlog.Error(ctx, "whatsAppMediator.persistMessageInElastic", fmt.Sprint("error when persist message", err))
		return err
	}

	return nil
}

func (w *whatsAppMediator) ReceiveMessage(ctx context.Context, received *dto.WhatsAppJSONReceived) error {

	chat := entity.Chat{
		Type:        entity.ChatText,
		MessageId:   received.Metadata.MessageId,
		OwnerId:     received.Owner.PhoneNumberId,
		ToPhoneId:   received.Sender.PhoneNumberId,
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
			Message:     entityChat.Message,
			ToPhone:     entityChat.ToPhoneId,
			ProfileName: entityChat.ProfileName,
		}); err != nil {
			hlog.Error(ctx, "whatsAppMediator.ReceiveMessage", fmt.Sprint("error when persist message", err))
		}
	}(&chat)

	return nil
}

func (w *whatsAppMediator) sendMessageToWhatsApp(ctx context.Context, ownerId string, message *dto.GatewayWhatsAppMessageDto) (string, error) {
	resp, err := w.whatsAppGateway.SendMessage(ctx, ownerId, message)

	if err != nil {
		hlog.Error(ctx, "whatsAppMediator.sendMessageToWhatsApp", fmt.Sprint("error when send message", err))
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

func (w *whatsAppMediator) createTextMessage(ctx context.Context, to, body string) *dto.GatewayWhatsAppMessageDto {
	return &dto.GatewayWhatsAppMessageDto{
		MessagingProduct: w.messageProduct,
		To:               to,
		Type:             dto.TextMessageType,
		// RecipientType:    "individual",
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
		Template: &dto.TemplateBody{
			Name: name,
			Language: dto.Language{
				Code: "pt-br",
			},
		},
	}
}
