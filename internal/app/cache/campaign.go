//go:generate mockgen -source=campaign.go -destination=mocks/campaign_mock.go -package=mocks

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/inter-hubly/pilot/database/hredis"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type Campaign interface {
	GetCampaignByPhoneId(ctx context.Context, campaign string) (*entity.Campaign, error)
	SaveCampaign(ctx context.Context, phoneId string, campaign *entity.Campaign) error
	GetStepInCampaign(ctx context.Context, phoneId string, index uint8) (*entity.Flow, error)
	GetIaContextInCampaign(ctx context.Context, phoneId string) (string, error)
}

type campaignCache struct {
	connection hredis.RedisConn
	database   string
}

var (
	_campaignCacheOnce sync.Once
	_campaignCache     *campaignCache
)

func NewCampaign(ctx context.Context) *campaignCache {
	_campaignCacheOnce.Do(func() {
		_campaignCache = &campaignCache{
			connection: hredis.GetConnection(ctx),
		}
	})
	return _campaignCache
}

func (c *campaignCache) GetCampaignByPhoneId(ctx context.Context, phoneId string) (*entity.Campaign, error) {
	hlog.Debug(ctx, "campaignCache.GetCampaignById", fmt.Sprintf("get campaign with phoneId: %s", phoneId))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s-campaign", tenantId, phoneId)

	campaignString, err := c.connection.GetClient(ctx).Get(ctx, key).Result()
	if err != nil {
		hlog.Error(ctx, "campaignCache.GetCampaignById", fmt.Sprintf("campaignId: %s", err.Error()))
		return nil, err
	}
	var campaign *entity.Campaign

	if err = json.Unmarshal([]byte(campaignString), &campaign); err != nil {
		hlog.Error(ctx, "campaignCache.GetCampaignById", fmt.Sprintf("campaignId: %s", err.Error()))
		return nil, err
	}
	return campaign, nil
}

func (c *campaignCache) SaveCampaign(ctx context.Context, phoneId string, campaign *entity.Campaign) error {
	hlog.Debug(ctx, "campaignCache.SaveCampaign", fmt.Sprintf("campaignId: %s", campaign.Name))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s-campaign", tenantId, phoneId)

	marshal, err := json.Marshal(campaign)
	if err != nil {
		hlog.Error(ctx, "campaignCache.SaveCampaign", fmt.Sprintf("campaignId: %s", campaign.Name))
		return err
	}
	if err = c.connection.GetClient(ctx).Set(ctx, key, marshal, 0).Err(); err != nil {
		hlog.Error(ctx, "campaignCache.SaveCampaign", fmt.Sprintf("campaignId: %s", campaign.Name))
		return err
	}
	return nil
}

func (c *campaignCache) GetStepInCampaign(ctx context.Context, phoneId string, index uint8) (*entity.Flow, error) {
	hlog.Debug(ctx, "campaignCache.GetStepInCampaign", fmt.Sprintf("campaignId for phoneId: %s with Index %d", phoneId, index))
	campaign, err := c.GetCampaignByPhoneId(ctx, phoneId)
	if err != nil {
		hlog.Error(ctx, "campaignCache.GetStepInCampaign", fmt.Sprintf("campaignId for phoneId: %s with Index %d", phoneId, index))
		return nil, err
	}
	if value, ok := campaign.Flows[strconv.Itoa(int(index))]; ok {
		return value, nil
	}
	// não achou nem o step da campanha, porém não posso considerado um erro
	return nil, nil
}

func (c *campaignCache) GetIaContextInCampaign(ctx context.Context, phoneId string) (string, error) {
	hlog.Debug(ctx, "campaignCache.GetIaContextInCampaign", fmt.Sprintf("phoneId: %s", phoneId))
	campaign, err := c.GetCampaignByPhoneId(ctx, phoneId)
	if err != nil {
		hlog.Error(ctx, "campaignCache.GetIaContextInCampaign", fmt.Sprintf("phoneId: %s", phoneId))
		return "", err
	}
	return campaign.IaContext, nil
}
