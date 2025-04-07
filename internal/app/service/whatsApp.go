//go:generate mockgen -source=whatsApp.go -destination=mocks/whatsApp_mock.go -package=mocks

package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/cache"
	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/linker/internal/app/mediator"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/entity"
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
	flowContext        cache.FlowContext
	chatGptGateway     gateway.Chatgpt
	campaignCache      cache.Campaign
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
			flowContext:        cache.NewFlowContext(ctx),
			chatGptGateway:     gateway.NewChatgpt(ctx),
			campaignCache:      cache.NewCampaign(ctx),
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
		hlog.Error(ctx, "whatsAppService.ReceiveMessage", err.Error())
		// return err
	}

	flowCount, err := w.flowContext.GetFlowCount(ctx, receivedDto.Sender.PhoneNumberId)
	if err != nil {
		hlog.Error(ctx, "whatsAppService.ReceiveMessage", err.Error())
		return err
	}
	flowEntity, err := w.campaignCache.GetStepInCampaign(ctx, receivedDto.Sender.PhoneNumberId, flowCount+1)
	if err != nil {
		hlog.Error(ctx, "whatsAppService.ReceiveMessage", err.Error())
		return err
	}

	if !flowEntity.HasIaInteraction {
		if err = w.createFlowResponse(ctx, receivedDto, flowEntity); err != nil {
			hlog.Error(ctx, "whatsAppService.ReceiveMessage", err.Error())
			return err
		}
		if err = w.flowContext.IncrementFlowCount(ctx, receivedDto.Sender.PhoneNumberId); err != nil {
			hlog.Error(ctx, "whatsAppService.ReceiveMessage", err.Error())
			return err
		}
		return nil
	}

	flowContext, err := w.flowContext.GetContext(ctx, receivedDto.Sender.PhoneNumberId)
	if err != nil {
		return err
	}
	return w.createIaFlowResponse(ctx, receivedDto, flowContext)
}

func (w *whatsAppService) createFlowResponse(ctx context.Context,
	receivedDto *dto.WhatsAppJSONReceived,
	flowEntity *entity.Flow,
) error {
	hlog.Debug(ctx, "whatsAppService.ReceiveMessage.createFlowResponse", fmt.Sprintf("%v", receivedDto))

	if err := w.SendMessage(ctx, &base.SendTextDto{
		To:      receivedDto.Sender.PhoneNumberId,
		Message: flowEntity.Message,
		IsOwner: true,
	}); err != nil {
		hlog.Error(ctx, "whatsAppService.ReceiveMessage.createFlowResponse", err.Error())
		return err
	}
	return nil
}

func (w *whatsAppService) createIaFlowResponse(ctx context.Context,
	receivedDto *dto.WhatsAppJSONReceived,
	flowContext []entity.Flow,
) error {
	hlog.Debug(ctx, "whatsAppService.ReceiveMessage.createIaFlowResponse", fmt.Sprintf("%v", receivedDto))
	var (
		messageToSend string
		err           error
	)

	flowContextMessage := &entity.Flow{
		HasIaInteraction: true,
		Message:          receivedDto.Metadata.Body,
	}

	messageToSend, err = w.chatGptGateway.GetInformation(ctx, flowContextMessage, flowContext)
	if err != nil {
		return err
	}

	// se tiver um contexto maior, é interacao de IA
	// if iaInteraction {
	// preciso adicionar a mensagem da IA no contexto

	// } else {
	// flowContextMessage.Message = flowContext[0].Message
	// messageToSend = flowContext[0].Message
	// }

	if _, err = w.flowContext.SaveContext(ctx, receivedDto.Sender.PhoneNumberId, flowContextMessage); err != nil {
		return err
	}

	err = w.SendMessage(ctx, &base.SendTextDto{
		To:      receivedDto.Sender.PhoneNumberId,
		Message: messageToSend,
		IsOwner: true,
	})

	if _, err = w.flowContext.SaveContext(ctx, receivedDto.Sender.PhoneNumberId, &entity.Flow{
		Role:    "assistant",
		Message: messageToSend,
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
