package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	dto "github.com/inter-hubly/linker/internal/domain/dto/whatsapp"
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

		whatsApp.StartTemplate()
		whatsApp.SendMessage()
		whatsApp.ChangeStatus()
		whatsApp.ReceiveMessage()
	})
}

type WhatsApp interface {
	StartTemplate()
	SendMessage()
	ChangeStatus()
	ReceiveMessage()
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
		if startTemplate.Name == "" {
			hlog.Error("whatsAppController.StartTemplate", "Name can't be empty")
		}
		w.whatsAppService.StartTemplate(ctx, &startTemplate)
	})
}

// SendMessage is when user send message in app
func (w *whatsAppController) SendMessage() {
	w.rabbit.Consume("whatsapp.send", func(value amqp.Delivery) {
		ctx := context.Background()
		var sentText dto.SendTextDto
		if err := json.Unmarshal(value.Body, &sentText); err != nil {
			hlog.Error("whatsAppController.SendMessage", fmt.Sprintf("err parsing: %s", err))
			return
		}
		if sentText.Message == "" {
			hlog.Error("whatsAppController.SendMessage", "Message can't be empty")
			return
		}
		w.whatsAppService.SendMessage(ctx, &sentText)
	})
}

// ReceiveMessage is when other numbers send message to me
func (w *whatsAppController) ReceiveMessage() {
	w.rabbit.Consume("whatsapp.message", func(value amqp.Delivery) {
		ctx := context.Background()
		var changeStatusDto dto.WhatsAppJSONReceived
		if err := json.Unmarshal(value.Body, &changeStatusDto); err != nil {
			hlog.Error("whatsAppController.ReceiveMessage", fmt.Sprintf("err parsing: %s", err))
			return
		}
		w.whatsAppService.ReceiveMessage(ctx, &changeStatusDto)
	})
}

func (w *whatsAppController) ChangeStatus() {
	w.rabbit.Consume("whatsapp.statuses", func(value amqp.Delivery) {
		ctx := context.Background()
		var changeStatusDto dto.ChangeStatusDto
		if err := json.Unmarshal(value.Body, &changeStatusDto); err != nil {
			hlog.Error("whatsAppController.ChangeStatus", fmt.Sprintf("err parsing: %s", err))
			return
		}
		if changeStatusDto.MessageId == "" {
			hlog.Error("whatsAppController.ChangeStatus", "Message can't be empty")
			return
		}
		w.whatsAppService.ChangeStatusMessage(ctx, &changeStatusDto)
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
