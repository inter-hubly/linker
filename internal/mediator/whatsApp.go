package mediator

import (
	"context"
	"errors"
	"time"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/linker/internal/gateway"
	"github.com/inter-hubly/linker/internal/repository"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	SentMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error
	DeliveredMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error
	SetStatus(ctx context.Context, message *entity.WhatsAppJSONReceived) error
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
		To:               template.To,
		Type:             dto.TemplateMessageType,
		Template: &dto.TemplateDto{
			Name: template.Name,
			LanguageDto: dto.LanguageDto{
				Code: template.Language,
			},
		},
	}

	if err := w.whatsAppGateway.SendMessage(ctx, "", &message); err != nil {
		hlog.Error("whatsAppMediator.StartTemplate", "error when send message to whatsApp", err)
		return err
	}

	// Id        string          `json:"id"`
	// MessageId string          `json:"messageId"`
	// Type      ChatType        `json:"type"`
	// Received  ChatMessageTime `json:"received,omitempty"`
	// Read      ChatMessageTime `json:"read,omitempty"`
	// Delivered ChatMessageTime `json:"delivered,omitempty"`
	// Message   ReceivedMessage `json:"message,omitempty"`

	entity := entity.Chat{}
	w.whatsAppRepository.PersistMessage(ctx)

	return nil
}

func (w *whatsAppMediator) SentMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	chanError := make(chan *errValue)

	messageToWhats := dto.GatewayWhatsAppMessageDto{
		MessagingProduct: w.messageProduct,
		To:               message.SenderPhone,
		Type:             dto.TemplateMessageType,
		Text: &dto.WhatsAppTextDto{
			PreviewUrl: false,
			Body:       message.Metadata.Body,
		},
	}
	go w.sendMessageToWhatsApp(ctx, &messageToWhats, chanError)

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

func (w *whatsAppMediator) DeliveredMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	chanError := make(chan *errValue)
	go w.sendMessageToWhatsApp(ctx, message, chanError)

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

func (w *whatsAppMediator) SetStatus(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	w.whatsAppRepository.SetStatusMessageById(ctx, message.Metadata.MessageId, message.Status, entity.ChatMessageTime{
		CreatedInDatabase: time.Now(),
		ReceivedAt:        message.Metadata.Timestamp,
	})
	return nil
}

func (w *whatsAppMediator) persistMessageInElastic(ctx context.Context, message *entity.WhatsAppJSONReceived, chanError chan *errValue) {
	normalizedChat := entity.NormalizeWhatsAppMessage(message)
	_, err := w.whatsAppRepository.PersistMessage(ctx, normalizedChat)
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
	err := w.whatsAppGateway.SendMessage(ctx, "message", message)
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
