package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type Campaign interface {
	GetCampaignById(ctx context.Context, flowId string) (*entity.Campaign, error)
}

type campaignRepository struct {
}

func NewCampaign() *campaignRepository {
	var (
		once       sync.Once
		repository *campaignRepository
	)

	once.Do(func() {
		repository = &campaignRepository{}
	})
	return repository
}

func (r *campaignRepository) GetCampaignById(ctx context.Context, flowId string) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaignRepository.GetCampaignById", fmt.Sprint("flowId", flowId))
	return &entity.Campaign{
		ParametersLength: 2,
		TemplateLanguage: "pt_BR",
		TemplateName:     "cobranca_mensal",
		Flows: []string{
			"12345",
		},
		Phones: []string{
			"48991784586",
			"48988356622",
		},
	}, nil
}
