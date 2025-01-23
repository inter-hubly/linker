package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
)

type Campaign interface {
	StartCampaign(ctx context.Context, campaignId string) error
}
type campaignService struct {
	campaignRepository  repository.Campaign
	contactRepository   repository.Contact
	variablesRepository repository.Variables
	whatsAppService     WhatsApp
}

func NewCampaign(ctx context.Context) *campaignService {

	var (
		serviceOnce sync.Once
		service     *campaignService
	)

	serviceOnce.Do(func() {
		service = &campaignService{
			campaignRepository:  repository.NewCampaign(),
			contactRepository:   repository.NewContact(),
			variablesRepository: repository.NewVariables(ctx),
			whatsAppService:     NewWhatsApp(),
		}
	})
	return service
}

func (s *campaignService) StartCampaign(ctx context.Context, campaignId string) error {
	hlog.Debug(ctx, "campaignService.StartCampaign", fmt.Sprint("campaignId", campaignId))
	campaign, err := s.campaignRepository.GetCampaignById(ctx, campaignId)
	if err != nil {
		hlog.Error(ctx, "campaignService.StartCampaign", err.Error())
		return err
	}

	contacts, err := s.contactRepository.GetContactsById(ctx, campaign.ContactsId...)
	if err != nil {
		hlog.Error(ctx, "campaignService.StartCampaign", err.Error())
		return err
	}

	parameters := make([]string, 0, len(campaign.Parameters))
	for _, p := range campaign.Parameters {
		parameters = append(parameters, p.Value)
	}

	for _, contact := range contacts {
		var variables map[string]interface{}
		variables, err = s.variablesRepository.GetVariablesByUserId(ctx, contact.Id.String(), parameters...)
		if err != nil {
			hlog.Error(ctx, "campaignService.StartCampaign", err.Error())
			return err
		}
		userParameters := make([]valueobject.Pair[string, string], 0, len(variables))
		for _, p := range parameters {
			v := valueobject.Pair[string, string]{
				Key:   "text",
				Value: variables[p].(string),
			}
			userParameters = append(userParameters, v)
		}

		err = s.whatsAppService.StartTemplate(ctx, &base.StartTemplateDto{
			To:         contact.Phone,
			CampaignId: campaign.Id,
			Template: base.TemplateInfo{
				Name:     campaign.Template.Name,
				Language: campaign.Template.Language,
			},
			Parameters: userParameters,
		})
	}
	return nil
}
