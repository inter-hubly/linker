package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type Flow interface {
	Start(ctx context.Context, flowId string) (string, error)
}

type flowService struct {
	flowRepository repository.Flow
	chatgptGateway gateway.Chatgpt
}

var (
	_flowServiceOnce sync.Once
	_flowService     *flowService
)

func NewFlow(ctx context.Context) *flowService {
	_flowServiceOnce.Do(func() {
		_flowService = &flowService{
			flowRepository: repository.NewFlow(ctx),
			chatgptGateway: gateway.NewChatgpt(ctx),
		}
	})
	return _flowService
}

func (s *flowService) Start(ctx context.Context, flowId string) (string, error) {
	hlog.Debug(ctx, "flowService.Start", fmt.Sprint("Start flow ", flowId))

	flowEntity, err := s.flowRepository.GetFlowById(ctx, flowId)
	if err != nil {
		hlog.Error(ctx, "flowService.Start", fmt.Sprint("Flow ", flowId, " GetFlow err:", err))
		return "", err
	}

	if flowEntity.HasAiCheck {
		return s.chatgptGateway.GetInformation(ctx, flowId)
	}

	return flowEntity.Message, nil
}
