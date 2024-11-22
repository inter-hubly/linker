package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/service"
	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/domain/dto"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/streadway/amqp"
)

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppController
)

func NewWhatsApp() {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppController{
			rabbit:          broker.GetConnection(),
			whatsAppService: service.NewWhatsApp(),
		}
		whatsApp.SentMessage()
		whatsApp.ReceiveMessage()
	})
}

type WhatsApp interface {
	SentMessage()
	ReceiveMessage()
}

type whatsAppController struct {
	exchange        string
	rabbit          broker.Connection
	whatsAppService service.WhatsApp
}

func (w *whatsAppController) SentMessage() {
	w.rabbit.Consume("whatsapp.sent", func(value amqp.Delivery) {
		ctx := context.Background()
		receivedDto, err := parseJsonReceived(ctx, value.Body)
		if err != nil {
			hlog.Error("whatsAppController.ReceiveMessage", fmt.Sprintf("err parsing: %s", err))
		}
		w.whatsAppService.SendMessage(ctx, receivedDto)
	})
}

func (w *whatsAppController) ReceiveMessage() {
	w.rabbit.Consume("whatsapp.receive", func(value amqp.Delivery) {
		ctx := context.Background()
		receivedDto, err := parseJsonReceived(ctx, value.Body)
		if err != nil {
			hlog.Error("whatsAppController.ReceiveMessage", fmt.Sprintf("err parsing: %s", err))
		}
		w.whatsAppService.ReceiveMessage(ctx, receivedDto)
	})
}

func parseJsonReceived(_ context.Context, body []byte) (*dto.WhatsAppJSONReceived, error) {
	hlog.Debug("whatsAppController.parseJsonReceived", fmt.Sprintf("%s", body))
	var receivedDto dto.WhatsAppJSONReceived
	if err := json.Unmarshal(body, &receivedDto); err != nil {
		return nil, err
	}
	return &receivedDto, nil
}
