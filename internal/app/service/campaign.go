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
	campaignRepository repository.Campaign
	contactRepository  repository.Contact
	whatsAppService    WhatsApp
}

var (
	_campaignServiceOnce sync.Once
	_campaignService     *campaignService
)

func NewCampaign(ctx context.Context) *campaignService {
	_campaignServiceOnce.Do(func() {
		_campaignService = &campaignService{
			campaignRepository: repository.NewCampaign(ctx),
			contactRepository:  repository.NewContact(ctx),
			whatsAppService:    NewWhatsApp(ctx),
		}
	})
	return _campaignService
}

func (s *campaignService) StartCampaign(ctx context.Context, campaignId string) error {
	hlog.Debug(ctx, "campaignService.StartCampaign", fmt.Sprint("campaignId", campaignId))
	campaignDb, err := s.campaignRepository.GetCampaignById(ctx, campaignId)
	if err != nil {
		hlog.Error(ctx, "campaignService.StartCampaign", err.Error())
		return err
	}

	if campaignDb.HasIaInteraction {
		hlog.Debug(ctx, "campaignService.StartCampaign", "IaInteraction")

	}

	contacts, err := s.contactRepository.GetContactsById(ctx, campaignDb.ContactsId...)
	if err != nil {
		hlog.Error(ctx, "campaignService.StartCampaign", err.Error())
		return err
	}

	// percorrer cada contato
	for _, contact := range contacts {
		userParameters := make([]valueobject.Pair[string, string], 0, len(campaignDb.Variables))

		// preciso verificar se o contato tem a variável necessária
		// então vou percorrer as variaveis que a campanha pede
		for _, p := range campaignDb.Variables {
			var name string
			name, err = contact.GetVariableByName(p.Value)
			if err != nil {
				hlog.Error(ctx, "campaignService.StartCampaign", err.Error())
				return err
			}
			userParameters = append(userParameters, valueobject.Pair[string, string]{Key: p.Key, Value: name})
		}

		err = s.whatsAppService.StartTemplate(ctx, &base.StartTemplateDto{
			To:         contact.Phone,
			CampaignId: campaignDb.Id,
			Template: base.TemplateInfo{
				Name:     campaignDb.Template.Name,
				Language: campaignDb.Template.Language,
			},
			Parameters: userParameters,
		})
	}
	return nil
}
