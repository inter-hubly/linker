package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/hredis"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
)

type IaContext interface {
	SaveContext(ctx context.Context, senderId string) (string, error)
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

func (r *iaContextRepository) SaveContext(ctx context.Context, senderId string) (string, error) {
	hlog.Debug(ctx, "iaContextRepository.SaveContext", fmt.Sprint("save ia context ", senderId))
	tenantId := hctx.Tenant.Get(ctx)
	r.connection.GetClient(ctx).Set(ctx, fmt.Sprintf("%s-%s", tenantId, senderId), "", 0)
}
