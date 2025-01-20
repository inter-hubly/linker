package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/entity"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
)

type Flows interface {
	GetFlowById(ctx context.Context, flowId ...string) (*entity.Flow, error)
}

type flowRepository struct {
}

func NewFlow() *flowRepository {
	clients = make(map[string]*valueobject.Client)

	var (
		clientOnce sync.Once
		client     *flowRepository
	)

	clientOnce.Do(func() {
		client = &flowRepository{}
	})
	return client
}

func (r *flowRepository) GetFlowById(ctx context.Context, flowId ...string) (*entity.Flow, error) {
	hlog.Debug(ctx, "flowRepository.GetFlowsById", fmt.Sprint("flowIds", flowId))
	return &entity.Flow{
		FlowType:  entity.Messaging,
		NextSteps: "123",
	}, nil
}
