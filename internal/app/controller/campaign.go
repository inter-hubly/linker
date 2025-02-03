package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/service"
	"github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/streadway/amqp"
)

func NewFlow(ctx context.Context) {

	var (
		controllerOnce sync.Once
		controller     *campaignController
	)

	controllerOnce.Do(func() {
		controller = &campaignController{
			rabbit:          broker.GetConnection(),
			campaignService: service.NewCampaign(ctx),
		}
	})

	controller.Init(ctx)
}

type Campaign interface {
	Init(ctx context.Context)
}

type campaignController struct {
	exchange        string
	rabbit          broker.Connection
	campaignService service.Campaign
}

func (c *campaignController) Init(ctx context.Context) {
	c.rabbit.Consume(ctx, "campaign.init", func(value amqp.Delivery) {
		header := value.Headers["tenantId"].(string)
		ctx = hctx.Tenant.New(header)

		var campaignInitDto struct {
			Id string `json:"id"`
		}

		if err := json.Unmarshal(value.Body, &campaignInitDto); err != nil {
			hlog.Error(ctx, "campaignController.Init", fmt.Sprintf("err parsing: %s", err))
			return
		}

		c.campaignService.StartCampaign(ctx, campaignInitDto.Id)
	})
}
