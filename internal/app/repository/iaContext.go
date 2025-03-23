package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/pilot/database/hredis"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type IaContext interface {
	StartContext(ctx context.Context, senderId, iaCompanyContext string) error
	SaveContext(ctx context.Context, senderId string, iaCompanyContext *dto.IaContext) (string, error)
	GetContext(ctx context.Context, senderId string) ([]dto.IaContext, error)
}

type iaContextRepository struct {
	connection hredis.RedisConn
	database   uint8
}

var (
	_iaContextRepositoryOnce sync.Once
	_iaContextRepository     *iaContextRepository
)

func NewIaContext(ctx context.Context) *iaContextRepository {
	_iaContextRepositoryOnce.Do(func() {
		_iaContextRepository = &iaContextRepository{
			connection: hredis.GetConnection(ctx),
		}
	})
	return _iaContextRepository
}

func (r *iaContextRepository) StartContext(ctx context.Context, senderId, iaCompanyContext string) error {
	hlog.Debug(ctx, "iaContextRepository.StartContext", fmt.Sprintf("start ia context %s", senderId))
	if _, err := r.SaveContext(ctx, senderId, &dto.IaContext{
		Role:    "system",
		Content: iaCompanyContext,
	}); err != nil {
		return err
	}
	return nil
}

func (r *iaContextRepository) SaveContext(ctx context.Context, senderId string, iaCompanyContext *dto.IaContext) (string, error) {
	hlog.Debug(ctx, "iaContextRepository.SaveContext", fmt.Sprintf("save ia context %s", senderId))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s", tenantId, senderId)

	marshal, err := json.Marshal(iaCompanyContext)
	if err != nil {
		hlog.Error(ctx, "iaContextRepository.SaveContext", fmt.Sprintf("error when marshal %s", iaCompanyContext))
		return "", err
	}
	if err = r.connection.GetClient(ctx).LPush(ctx, key, marshal).Err(); err != nil {
		hlog.Error(ctx, "iaContextRepository.SaveContext", fmt.Sprintf("save ia context %s", senderId))
		return "", err
	}

	return tenantId, nil
}

func (r *iaContextRepository) GetContext(ctx context.Context, senderId string) ([]dto.IaContext, error) {
	hlog.Debug(ctx, "iaContextRepository.GetContext", fmt.Sprint("get ia context ", senderId))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s", tenantId, senderId)

	messages, err := r.connection.GetClient(ctx).LRange(ctx, key, 0, -1).Result()
	if err != nil {
		hlog.Error(ctx, "iaContextRepository.GetContext", fmt.Sprintf("get ia context %s", senderId))
		return nil, err
	}
	if len(messages) == 0 {
		return nil, errors.New("ia context not found")
	}
	contexts := make([]dto.IaContext, 0, len(messages))
	for _, message := range messages {
		var iaContext dto.IaContext
		if err = json.Unmarshal([]byte(message), &iaContext); err != nil {
			hlog.Error(ctx, "iaContextRepository.GetContext", fmt.Sprintf("unmarshal %s", message))
			return nil, err
		}
		contexts = append(contexts, iaContext)
	}
	return contexts, nil
}
