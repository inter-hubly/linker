package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/server"
)

type Flow interface {
	Start(ctx context.Context, iaMessage *entity.Flow, iaContext []entity.Flow) (string, error)
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
	var chatGptGateway gateway.Chatgpt
	if server.GetEnvironment().Env != "development" {
		chatGptGateway = gateway.NewChatGptMock()
	} else {
		chatGptGateway = gateway.NewChatgpt(ctx)
	}

	_flowServiceOnce.Do(func() {
		_flowService = &flowService{
			flowRepository: repository.NewFlow(ctx),
			chatgptGateway: chatGptGateway,
		}
	})
	return _flowService
}

func (s *flowService) Start(ctx context.Context, iaMessage *entity.Flow, iaContext []entity.Flow) (string, error) {
	hlog.Debug(ctx, "flowService.Start", fmt.Sprint("Start flow with context count", len(iaContext)))

	information, err := s.chatgptGateway.GetInformation(ctx, iaMessage, iaContext)
	if err != nil {
		hlog.Error(ctx, "flowService.Start", fmt.Sprint("Failed to get information from count ", len(iaContext)))
		return "", err
	}
	return information, nil
}
