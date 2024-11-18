package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/inter-hubly/linker/internal/domain/dto"
	"sync"

	"github.com/inter-hubly/linker/internal/service"
	"github.com/inter-hubly/pilot/broker"
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
		whatsApp.SendMessage()
		whatsApp.ReceiveMessage()
	})
}

type WhatsApp interface {
	SendMessage()
	ReceiveMessage()
}

type whatsAppController struct {
	exchange        string
	rabbit          broker.Connection
	whatsAppService service.WhatsApp
}

func (w *whatsAppController) SendMessage() {
	w.rabbit.Consume("whatsapp.send", func(value amqp.Delivery) {
		ctx := context.Background()
		hlog.Debug("whatsAppController.SendMessage", fmt.Sprintf("%s", value.Body))
		var receivedDto dto.WhatsAppJSONReceived
		if err := json.Unmarshal(value.Body, &receivedDto); err != nil {
			hlog.Error("whatsAppController.SendMessage", fmt.Sprintf("err parsing: %s", err))
			return
		}
		w.whatsAppService.SendMessage(ctx, &receivedDto)
	})
}

func (w *whatsAppController) ReceiveMessage() {
	w.rabbit.Consume("whatsapp.receive", func(value amqp.Delivery) {
		ctx := context.Background()
		hlog.Debug("whatsAppController.ReceiveMessage", fmt.Sprintf("%s", value.Body))
		var receivedDto dto.WhatsAppJSONReceived
		if err := json.Unmarshal(value.Body, &receivedDto); err != nil {
			hlog.Error("whatsAppController.ReceiveMessage", fmt.Sprintf("err parsing: %s", err))
			return
		}

		w.whatsAppService.ReceiveMessage(ctx, &receivedDto)
	})
}
