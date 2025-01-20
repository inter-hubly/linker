package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type Campaign interface {
	StartCampaign(ctx context.Context, campaignId string) error
}
type campaignService struct {
	campaignRepository repository.Campaign
	flowsRepository    repository.Flows
	whatsAppService    WhatsApp
}

func NewCampaign() *campaignService {

	var (
		serviceOnce sync.Once
		service     *campaignService
	)

	serviceOnce.Do(func() {
		service = &campaignService{
			campaignRepository: repository.NewCampaign(),
			flowsRepository:    repository.NewFlow(),
			whatsAppService:    NewWhatsApp(),
		}
	})
	return service
}

func (s *campaignService) StartCampaign(ctx context.Context, campaignId string) error {
	hlog.Debug(ctx, "campaignService.StartCampaign", fmt.Sprint("campaignId", campaignId))
	// campaign, err := s.campaignRepository.GetCampaignById(ctx, campaignId)
	// if err != nil {
	// 	hlog.Error(ctx, "campaignService.StartCampaign", err.Error())
	// 	return err
	// }

	// flow, err := s.flowsRepository.GetFlowById(ctx, campaign.Flows...)
	// if err != nil {
	// 	hlog.Error(ctx, "campaignService.StartCampaign", err.Error())
	// 	return err
	// }
	//
	// for _, v := range campaign.Phones {
	// 	s.whatsAppService.StartTemplate(ctx, &dto.StartTemplateDto{})
	// }
	return nil
}
