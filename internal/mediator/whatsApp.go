package mediator

import (
	"context"
	"errors"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/linker/internal/gateway"
	"github.com/inter-hubly/linker/internal/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	Persist(ctx context.Context, message *dto.WhatsAppJSONReceived) error
}

type whatsAppMediator struct {
	whatsAppRepository repository.WhatsApp
	whatsAppGateway    gateway.WhatsApp
}

func NewWhatsApp() WhatsApp {
	return &whatsAppMediator{}

}

func (w *whatsAppMediator) Persist(ctx context.Context, message *dto.WhatsAppJSONReceived) error {

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

func (w *whatsAppMediator) persistMessageInElastic(ctx context.Context, message *dto.WhatsAppJSONReceived, chanError chan *errValue) {
	func() {
		err := w.whatsAppRepository.PersistMessage(ctx, message)
		if err != nil {
			hlog.Error("whatsAppMediator.persistMessageInElastic", "error when persist message", err)
			chanError <- &errValue{
				errType: PersistError,
				err:     err,
			}
		}
		chanError <- nil
	}()
}

func (w *whatsAppMediator) sendMessageToWhatsApp(ctx context.Context, message *dto.WhatsAppJSONReceived, chanError chan *errValue) {
	func() {
		err := w.whatsAppGateway.SendMessage(ctx, message)
		if err != nil {
			hlog.Error("whatsAppMediator.sendMessageToWhatsApp", "error when send message", err)
			chanError <- &errValue{
				errType: SendError,
				err:     err,
			}
		}
		chanError <- nil
	}()
}

const (
	SendError    = "SendError"
	PersistError = "PersistError"
)

type errValue struct {
	errType string
	err     error
}
