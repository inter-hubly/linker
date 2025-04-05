package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/linker/internal/app/gateway"
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
	iaContext          repository.IaContext
	chatGptGateway     gateway.Chatgpt
}

var (
	_whatsAppServiceOnce sync.Once
	_whatsAppService     *whatsAppService
)

func NewWhatsApp(ctx context.Context) *whatsAppService {
	_whatsAppServiceOnce.Do(func() {
		_whatsAppService = &whatsAppService{
			whatsappMediator:   mediator.NewWhatsApp(ctx),
			whatsappRepository: repository.NewWhatsApp(ctx),
			campaignRepository: repository.NewCampaign(ctx),
			iaContext:          repository.NewIaContext(ctx),
			chatGptGateway:     gateway.NewChatgpt(ctx),
		}
	})
	return _whatsAppService
}

func (w *whatsAppService) StartTemplate(ctx context.Context, template *base.StartTemplateDto) error {
	hlog.Debug(ctx, "whatsAppService.StartTemplate", fmt.Sprintf("%v", template))

	components := w.createParameters(ctx, template.Parameters)
	gatewayTemplate := dto.GatewayWhatsAppMessageDto{
		MessagingProduct: "whatsapp",
		To:               template.To,
		Type:             dto.TemplateMessageType,
		Template: &dto.TemplateBody{
			Name: template.Template.Name,
			Language: dto.Language{
				Code: template.Template.Language,
			},
			Components: []dto.Component{
				{
					Type:       "body",
					Parameters: components,
				},
			},
		},
	}

	return w.whatsappMediator.StartTemplate(ctx, template.CampaignId, &gatewayTemplate)
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

func (w *whatsAppService) ReceiveMessage(ctx context.Context, receivedDto *dto.WhatsAppJSONReceived) error {
	hlog.Debug(ctx, "whatsAppService.ReceiveMessage", fmt.Sprintf("%v", receivedDto))

	if err := w.whatsappMediator.ReceiveMessage(ctx, receivedDto); err != nil {
		return err
	}

	iaContext, err := w.iaContext.GetContext(ctx, receivedDto.Sender.PhoneNumberId)
	if err == nil {
		return w.createAiResponse(ctx, receivedDto, iaContext)
	}

	return nil
}

func (w *whatsAppService) createAiResponse(ctx context.Context, receivedDto *dto.WhatsAppJSONReceived, iaContext []dto.IaContext) error {
	var (
		messageToSend string
		err           error
	)

	iaContextMessage := &dto.IaContext{
		Role:    "user",
		Content: receivedDto.Metadata.Body,
	}

	if _, err = w.iaContext.SaveContext(ctx, receivedDto.Sender.PhoneNumberId, iaContextMessage); err != nil {
		return err
	}

	messageToSend, err = w.chatGptGateway.GetInformation(ctx, iaContextMessage, iaContext)
	if err != nil {
		return err
	}

	err = w.SendMessage(ctx, &base.SendTextDto{
		To:      receivedDto.Sender.PhoneNumberId,
		Message: messageToSend,
		IsOwner: true,
	})

	if _, err = w.iaContext.SaveContext(ctx, receivedDto.Sender.PhoneNumberId, &dto.IaContext{
		Role:    "assistant",
		Content: messageToSend,
	}); err != nil {
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
