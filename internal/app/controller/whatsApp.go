package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/linker/internal/app/service"
	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/streadway/amqp"
)

func NewWhatsApp(ctx context.Context) {
	var (
		whatsAppOnce sync.Once
		whatsApp     *whatsAppController
	)

	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppController{
			rabbit:          broker.GetConnection(),
			whatsAppService: service.NewWhatsApp(),
		}

		whatsApp.StartTemplate(ctx)
		whatsApp.SendMessage(ctx)
		whatsApp.ChangeStatus(ctx)
		whatsApp.ReceiveMessage(ctx)
	})
}

type WhatsApp interface {
	StartTemplate(ctx context.Context)
	SendMessage(ctx context.Context)
	ChangeStatus(ctx context.Context)
	ReceiveMessage(ctx context.Context)
}

type whatsAppController struct {
	exchange        string
	rabbit          broker.Connection
	whatsAppService service.WhatsApp
}

func (w *whatsAppController) StartTemplate(ctx context.Context) {
	w.rabbit.Consume(ctx, "whatsapp.start", func(value amqp.Delivery) {
		header := value.Headers["tenantId"].(string)
		ctx = hctx.Tenant.New(header)

		var startTemplate base.StartTemplateDto
		if err := json.Unmarshal(value.Body, &startTemplate); err != nil {
			hlog.Error(ctx, "whatsAppController.StartTemplate", fmt.Sprintf("err parsing: %s", err))
			return
		}
		if startTemplate.CampaignId.String() == "" {
			hlog.Error(ctx, "whatsAppController.StartTemplate", "CampaignId can't be empty")
		}
		w.whatsAppService.StartTemplate(ctx, &startTemplate)
	})
}

// SendMessage is when user send message in app
func (w *whatsAppController) SendMessage(ctx context.Context) {
	w.rabbit.Consume(ctx, "whatsapp.send", func(value amqp.Delivery) {
		header := value.Headers["tenantId"].(string)
		ctx = hctx.Tenant.New(header)

		var sentText base.SendTextDto
		if err := json.Unmarshal(value.Body, &sentText); err != nil {
			hlog.Error(ctx, "whatsAppController.SendMessage", fmt.Sprintf("err parsing: %s", err))
			return
		}
		if sentText.Message == "" {
			hlog.Error(ctx, "whatsAppController.SendMessage", "Message can't be empty")
			return
		}
		if sentText.To == "" {
			hlog.Error(ctx, "whatsAppController.SendMessage", "Sender and Receiver can't be empty")
			return
		}
		w.whatsAppService.SendMessage(ctx, &sentText)
	})
}

// ReceiveMessage is when other numbers send message to me
func (w *whatsAppController) ReceiveMessage(ctx context.Context) {
	w.rabbit.Consume(ctx, "whatsapp.message", func(value amqp.Delivery) {
		header := value.Headers["tenantId"].(string)
		ctx = hctx.Tenant.New(header)

		var changeStatusDto dto.WhatsAppJSONReceived
		if err := json.Unmarshal(value.Body, &changeStatusDto); err != nil {
			hlog.Error(ctx, "whatsAppController.ReceiveMessage", fmt.Sprintf("err parsing: %s", err))
			return
		}
		w.whatsAppService.ReceiveMessage(ctx, &changeStatusDto)
	})
}

func (w *whatsAppController) ChangeStatus(ctx context.Context) {
	w.rabbit.Consume(ctx, "whatsapp.statuses", func(value amqp.Delivery) {
		header := value.Headers["tenantId"].(string)
		ctx = hctx.Tenant.New(header)

		var changeStatusDto dto.ChangeStatusDto
		if err := json.Unmarshal(value.Body, &changeStatusDto); err != nil {
			hlog.Error(ctx, "whatsAppController.ChangeStatus", fmt.Sprintf("err parsing: %s", err))
			return
		}
		if changeStatusDto.MessageId == "" {
			hlog.Error(ctx, "whatsAppController.ChangeStatus", "Message can't be empty")
			return
		}
		w.whatsAppService.ChangeStatusMessage(ctx, &changeStatusDto)
	})
}

func parseJsonReceived(ctx context.Context, body []byte) (*dto.WhatsAppJSONReceived, error) {
	hlog.Debug(ctx, "whatsAppController.parseJsonReceived", fmt.Sprintf("%s", body))
	var receivedDto dto.WhatsAppJSONReceived
	if err := json.Unmarshal(body, &receivedDto); err != nil {
		return nil, err
	}
	return &receivedDto, nil
}
