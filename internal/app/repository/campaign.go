package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/entity"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/valueobject"
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

func (r *campaignRepository) GetCampaignById(ctx context.Context, campaignId string) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaignRepository.GetCampaignById", fmt.Sprint("campaignId", campaignId))
	return &entity.Campaign{
		Parameters: []valueobject.Pair[string, string]{
			{"text", "name"},
			{"text", "value"},
		},
		Template: base.TemplateInfo{
			Language: "pt_BR",
			Name:     "cobranca_mensal",
		},
		ContactsId: []string{
			"a92ce761-d1eb-423a-9563-d2de43e888b2",
			"34927c6b-8dae-4c04-b2e1-5f21032f31f1",
		},
	}, nil
}
