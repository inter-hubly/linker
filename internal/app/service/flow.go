package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"
)

type Flow interface {
	Start(ctx context.Context, iaContext []dto.IaContext) (string, error)
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

func (s *flowService) Start(ctx context.Context, iaContext []dto.IaContext) (string, error) {
	hlog.Debug(ctx, "flowService.Start", fmt.Sprint("Start flow with context count", len(iaContext)))

	information, err := s.chatgptGateway.GetInformation(ctx, iaContext)
	if err != nil {
		hlog.Error(ctx, "flowService.Start", fmt.Sprint("Failed to get information from count ", len(iaContext)))
		return "", err
	}
	return information, nil
}
