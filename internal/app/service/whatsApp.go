package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/linker/internal/app/mediator"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	ChangeStatusMessage(ctx context.Context, message *dto.ChangeStatusDto) error
	SendMessage(ctx context.Context, template *base.SendTextDto) error
	StartTemplate(ctx context.Context, template *base.StartTemplateDto) error
	ReceiveMessage(ctx context.Context, dto *dto.WhatsAppJSONReceived) error
}

type whatsAppService struct {
	whatsappMediator   mediator.WhatsApp
	whatsappRepository repository.WhatsApp
	campaignRepository repository.Campaign
}

func NewWhatsApp() *whatsAppService {

	var (
		serviceOnce sync.Once
		whatsApp    *whatsAppService
	)

	serviceOnce.Do(func() {
		whatsApp = &whatsAppService{
			whatsappMediator:   mediator.NewWhatsApp(),
			whatsappRepository: repository.NewWhatsApp(),
			campaignRepository: repository.NewCampaign(),
		}
	})
	return whatsApp
}

func (w *whatsAppService) StartTemplate(ctx context.Context, template *base.StartTemplateDto) error {
	hlog.Debug(ctx, "whatsAppService.StartTemplate", fmt.Sprintf("%v", template))
	campaign, err := w.campaignRepository.GetCampaignById(ctx, template.CampaignId)
	if err != nil {
		hlog.Error(ctx, "whatsAppMediator.StartTemplate", err.Error())
		return err
	}
	if uint8(len(template.Parameters)) != campaign.ParametersLength {
		hlog.Error(ctx, "whatsAppMediator.StartTemplate", "wrong number of parameters")
	}

	components := w.createParameters(ctx, template.Parameters)
	gatewayTemplate := dto.GatewayWhatsAppMessageDto{
		MessagingProduct: "whatsapp",
		To:               template.To,
		Type:             dto.TemplateMessageType,
		Template: &dto.TemplateBody{
			Name: campaign.TemplateName,
			Language: dto.Language{
				Code: campaign.TemplateLanguage,
			},
			Components: []dto.Component{
				{
					Type:       "body",
					Parameters: components,
				},
			},
		},
	}

	return w.whatsappMediator.StartTemplate(ctx, campaign.Id, &gatewayTemplate)
}

func (w *whatsAppService) SendMessage(ctx context.Context, template *base.SendTextDto) error {
	hlog.Debug(ctx, "whatsAppService.SendMessage", fmt.Sprintf("%v", template))

	tenantId := hctx.Tenant.Get(ctx)
	message := dto.WhatsAppJSONReceived{
		Owner: dto.WhatsAppPhoneIdDto{
			PhoneNumberId: tenantId,
		},
		Sender: dto.WhatsAppPhoneIdDto{
			PhoneNumberId: template.To,
		},
		Metadata: dto.WhatsAppMetadataDto{
			Body: template.Message,
		},
	}

	return w.whatsappMediator.SendMessage(ctx, &message)
}

func (w *whatsAppService) ReceiveMessage(ctx context.Context, dto *dto.WhatsAppJSONReceived) error {
	hlog.Debug(ctx, "whatsAppService.ReceiveMessage", fmt.Sprintf("%v", dto))

	if err := w.whatsappMediator.ReceiveMessage(ctx, dto); err != nil {
		return err
	}

	return nil
}

func (w *whatsAppService) ChangeStatusMessage(ctx context.Context, message *dto.ChangeStatusDto) error {
	hlog.Debug(ctx, "whatsAppService.ChangeStatusMessage", fmt.Sprintf("%v", message))

	if err := w.whatsappRepository.SetStatusMessageById(ctx, message.MessageId, message.Status, message.ExpirationTimeStamp); err != nil {
		return err
	}
	return nil
}

func (w *whatsAppService) createParameters(ctx context.Context, components []valueobject.Pair[string, string]) []dto.Parameter {
	resp := make([]dto.Parameter, 0, len(components))
	for _, component := range components {
		resp = append(resp, dto.Parameter{
			Type: component.Key,
			Text: component.Value,
		})
	}
	return resp
}
