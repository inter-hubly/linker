//go:generate mockgen -source=flowContext.go -destination=mocks/flowContext_mock.go -package=mocks

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/inter-hubly/pilot/database/hredis"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type FlowContext interface {
	StartContext(ctx context.Context, senderId string, flowContext *entity.Flow) error
	SaveContext(ctx context.Context, senderId string, iaCompanyContext *entity.Flow) (string, error)
	GetContext(ctx context.Context, senderId string) ([]entity.Flow, error)
	GetFlowCount(ctx context.Context, senderId string) (uint8, error)
	IncrementFlowCount(ctx context.Context, senderId string) error
}

type flowContextRepository struct {
	connection hredis.RedisConn
	database   uint8
}

var (
	_flowContextRepositoryOnce sync.Once
	_flowContextRepository     *flowContextRepository
)

func NewFlowContext(ctx context.Context) *flowContextRepository {
	_flowContextRepositoryOnce.Do(func() {
		_flowContextRepository = &flowContextRepository{
			connection: hredis.GetConnection(ctx),
		}
	})
	return _flowContextRepository
}

func (r *flowContextRepository) StartContext(ctx context.Context, senderId string, campaign *entity.Flow) error {
	hlog.Debug(ctx, "flowContextRepository.StartContext", fmt.Sprintf("start ia context %s", senderId))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s-context", tenantId, senderId)

	marshal, err := json.Marshal(campaign)
	if err != nil {
		hlog.Error(ctx, "flowContextRepository.StartContext", fmt.Sprintf("error when marshal %v", campaign))
		return err
	}
	if err = r.connection.GetClient(ctx).LPush(ctx, key, marshal).Err(); err != nil {
		hlog.Error(ctx, "flowContextRepository.StartContext", fmt.Sprintf("save flow context %s", senderId))
		return err
	}

	return nil
}

func (r *flowContextRepository) SaveContext(ctx context.Context, senderId string, flowContext *entity.Flow) (string, error) {
	hlog.Debug(ctx, "flowContextRepository.SaveContext", fmt.Sprintf("save flow context %s", senderId))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s-context", tenantId, senderId)

	marshal, err := json.Marshal(flowContext)
	if err != nil {
		hlog.Error(ctx, "flowContextRepository.SaveContext", fmt.Sprintf("error when marshal %v", flowContext))
		return "", err
	}
	if !flowContext.IsIaInteraction {
		if err = r.connection.GetClient(ctx).LPop(ctx, key).Err(); err != nil {
			hlog.Error(ctx, "flowContextRepository.SaveContext", fmt.Sprintf("remove flow context %s", senderId))
			return "", err
		}
	}
	if err = r.connection.GetClient(ctx).LPush(ctx, key, marshal).Err(); err != nil {
		hlog.Error(ctx, "flowContextRepository.SaveContext", fmt.Sprintf("save flow context %s", senderId))
		return "", err
	}

	return tenantId, nil
}

func (r *flowContextRepository) GetContext(ctx context.Context, senderId string) ([]entity.Flow, error) {
	hlog.Debug(ctx, "flowContextRepository.GetContext", fmt.Sprint("get flow context ", senderId))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s-context", tenantId, senderId)

	messages, err := r.connection.GetClient(ctx).LRange(ctx, key, 0, -1).Result()
	if err != nil {
		hlog.Error(ctx, "flowContextRepository.GetContext", fmt.Sprintf("get flow context %s", senderId))
		return nil, err
	}
	if len(messages) == 0 {
		return nil, errors.New("flow context not found")
	}
	var contexts []entity.Flow
	messages, err = r.connection.GetClient(ctx).LRange(ctx, key, 0, -1).Result()

	for _, message := range messages {
		var entFlow entity.Flow
		if err = json.Unmarshal([]byte(message), &entFlow); err != nil {
			hlog.Error(ctx, "flowContextRepository.GetContext", fmt.Sprintf("unmarshal %s", message))
			return nil, err
		}
		contexts = append(contexts, entFlow)
	}

	return contexts, nil
}

func (r *flowContextRepository) GetFlowCount(ctx context.Context, senderId string) (uint8, error) {
	hlog.Debug(ctx, "flowContextRepository.GetFlowCount", fmt.Sprint("get flow context ", senderId))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s", tenantId, senderId)
	count, err := r.connection.GetClient(ctx).Get(ctx, key).Result()
	if err != nil {
		hlog.Error(ctx, "flowContextRepository.GetFlowCount", fmt.Sprintf("get flow context %s", key))
		if err = r.connection.GetClient(ctx).Set(ctx, key, "0", 0).Err(); err != nil {
			return 0, err
		}
		return 0, nil
	}
	val, err := strconv.Atoi(count)
	if err != nil {
		return 0, err
	}

	return uint8(val), nil
}

func (r *flowContextRepository) IncrementFlowCount(ctx context.Context, senderId string) error {
	hlog.Debug(ctx, "flowContextRepository.IncrementFlowCount", fmt.Sprint("get flow context ", senderId))
	tenantId := hctx.Tenant.Get(ctx)
	key := fmt.Sprintf("%s-%s", tenantId, senderId)
	_, err := r.connection.GetClient(ctx).Incr(ctx, key).Result()
	if err != nil {
		hlog.Error(ctx, "flowContextRepository.IncrementFlowCount", fmt.Sprintf("get flow context %s", senderId))
		return err
	}
	return nil
}
