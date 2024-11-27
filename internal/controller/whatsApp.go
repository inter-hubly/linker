package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/linker/internal/service"
	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/domain/entity"
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

		whatsApp.StartTemplate()
		whatsApp.DeliveredMessage()
		whatsApp.SentMessage()
		whatsApp.ReceivedMessage()
		whatsApp.ReadMessage()
	})
}

type WhatsApp interface {
	StartTemplate()
	SentMessage()
	ReceivedMessage()
	ReadMessage()
	DeliveredMessage()
}

type whatsAppController struct {
	exchange        string
	rabbit          broker.Connection
	whatsAppService service.WhatsApp
}

func (w *whatsAppController) StartTemplate() {
	w.rabbit.Consume("whatsapp.start", func(value amqp.Delivery) {
		ctx := context.Background()
		var startTemplate dto.StartTemplateDto
		if err := json.Unmarshal(value.Body, &startTemplate); err != nil {
			hlog.Error("whatsAppController.StartTemplate", fmt.Sprintf("err parsing: %s", err))
			return
		}
		w.whatsAppService.StartTemplate(ctx, &startTemplate)
	})
}

func (w *whatsAppController) SentMessage() {
	w.rabbit.Consume("whatsapp.sent", func(value amqp.Delivery) {
		ctx := context.Background()
		var sentText dto.SentTextDto
		if err := json.Unmarshal(value.Body, &sentText); err != nil {
			hlog.Error("whatsAppController.SentMessage", fmt.Sprintf("err parsing: %s", err))
			return
		}
		w.whatsAppService.SentMessage(ctx, &sentText)
	})
}

func (w *whatsAppController) DeliveredMessage() {
	w.rabbit.Consume("whatsapp.delivered", func(value amqp.Delivery) {
		ctx := context.Background()
		receivedDto, err := parseJsonReceived(ctx, value.Body)
		if err != nil {
			hlog.Error("whatsAppController.DeliveredMessage", fmt.Sprintf("err parsing: %s", err))
		}
		w.whatsAppService.DeliveredMessage(ctx, receivedDto)
	})
}

func (w *whatsAppController) ReceivedMessage() {
	w.rabbit.Consume("whatsapp.received", func(value amqp.Delivery) {
		ctx := context.Background()
		receivedDto, err := parseJsonReceived(ctx, value.Body)
		if err != nil {
			hlog.Error("whatsAppController.ReceivedMessage", fmt.Sprintf("err parsing: %s", err))
		}
		w.whatsAppService.SetMessageStatus(ctx, receivedDto)
	})
}

func (w *whatsAppController) ReadMessage() {
	w.rabbit.Consume("whatsapp.read", func(value amqp.Delivery) {
		ctx := context.Background()
		receivedDto, err := parseJsonReceived(ctx, value.Body)
		if err != nil {
			hlog.Error("whatsAppController.ReadMessage", fmt.Sprintf("err parsing: %s", err))
		}
		w.whatsAppService.SetMessageStatus(ctx, receivedDto)
	})
}

func parseJsonReceived(_ context.Context, body []byte) (*entity.WhatsAppJSONReceived, error) {
	hlog.Debug("whatsAppController.parseJsonReceived", fmt.Sprintf("%s", body))
	var receivedDto entity.WhatsAppJSONReceived
	if err := json.Unmarshal(body, &receivedDto); err != nil {
		return nil, err
	}
	return &receivedDto, nil
}
