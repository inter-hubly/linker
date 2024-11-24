package mediator

import (
	"context"
	"errors"
	"time"

	"github.com/inter-hubly/linker/internal/domain/entity"
	"github.com/inter-hubly/linker/internal/gateway"
	"github.com/inter-hubly/linker/internal/repository"
	"github.com/inter-hubly/pilot/domain/dto"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	SentMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
	DeliveredMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
	SetStatus(ctx context.Context, message *dto.WhatsAppJSONReceived) error
}

type whatsAppMediator struct {
	whatsAppRepository repository.WhatsApp
	whatsAppGateway    gateway.WhatsApp
}

func NewWhatsApp() WhatsApp {
	return &whatsAppMediator{
		whatsAppRepository: repository.NewWhatsApp(),
		whatsAppGateway:    gateway.NewWhatsApp(),
	}

}

func (w *whatsAppMediator) SentMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {

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

func (w *whatsAppMediator) DeliveredMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {
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

func (w *whatsAppMediator) SetStatus(ctx context.Context, message *dto.WhatsAppJSONReceived) error {
	w.whatsAppRepository.SetStatusMessageById(ctx, message.Metadata.MessageId, message.Status, entity.ChatMessageTime{
		CreatedInDatabase: time.Now(),
		ReceivedAt:        message.Metadata.Timestamp,
	})
	return nil
}

func (w *whatsAppMediator) persistMessageInElastic(ctx context.Context, message *dto.WhatsAppJSONReceived, chanError chan *errValue) {
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

func (w *whatsAppMediator) sendMessageToWhatsApp(ctx context.Context, message *dto.WhatsAppJSONReceived, chanError chan *errValue) {
	err := w.whatsAppGateway.SendMessage(ctx, message)
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
