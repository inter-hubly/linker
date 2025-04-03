//go:generate mockgen -source=campaign.go -destination=mocks/campaign_mock.go -package=mocks

package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Campaign interface {
	GetCampaignById(ctx context.Context, flowId string) (*entity.Campaign, error)
}

type campaignRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

var (
	campaignRepositoryOnce sync.Once
	campaign               *campaignRepository
)

func NewCampaign(ctx context.Context) *campaignRepository {

	campaignRepositoryOnce.Do(func() {
		campaign = &campaignRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "campaign",
		}
	})
	return campaign
}

func (r *campaignRepository) GetCampaignById(ctx context.Context, campaignId string) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaignRepository.GetCampaignById", fmt.Sprintf("campaignId: %s", campaignId))
	var camp entity.Campaign

	objID, err := primitive.ObjectIDFromHex(campaignId)
	if err != nil {
		hlog.Error(ctx, "campaignRepository.GetCampaignById", fmt.Sprintf("campaignId: %s", campaignId))
		return nil, err
	}

	if err = r.connection.GetCollection(ctx, r.collection).FindOne(ctx,
		bson.M{
			"_id": objID,
		},
	).Decode(&camp); err != nil {
		hlog.Error(ctx, "campaignRepository.GetCampaignById", fmt.Sprintf("campaignId: %s", campaignId))
		return nil, err
	}
	return &camp, nil
}
