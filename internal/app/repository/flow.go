package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/lentity"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hlog"
)

type Flow interface {
	GetFlowById(ctx context.Context, userId string) (*lentity.Flow, error)
}

type flowRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

var (
	_flowRepositoryOnce sync.Once
	_flowRepository     *flowRepository
)

func NewFlow(ctx context.Context) *flowRepository {
	_flowRepositoryOnce.Do(func() {
		_flowRepository = &flowRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "flow",
		}
	})
	return _flowRepository
}

func (r *flowRepository) GetFlowById(ctx context.Context, flowId string) (*lentity.Flow, error) {
	hlog.Debug(ctx, "flowRepository.GetFlowById", fmt.Sprint("Get flow by id :", flowId))
	return &lentity.Flow{
		Id:         "123456",
		Name:       "chatgpt",
		HasAiCheck: true,
	}, nil
}
