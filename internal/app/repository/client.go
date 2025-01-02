package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/hlog"
)

var clients map[string]*valueobject.Client

type Client interface {
	GetClientById(ctx context.Context, clientId string) (*valueobject.Client, error)
}

type clientRepository struct {
	keeperGateway gateway.Keeper
}

func NewClient() *clientRepository {
	clients = make(map[string]*valueobject.Client)

	var (
		clientOnce sync.Once
		client     *clientRepository
	)

	clientOnce.Do(func() {
		client = &clientRepository{
			keeperGateway: gateway.NewKeeper(),
		}
	})
	return client
}

func (r *clientRepository) GetClientById(ctx context.Context, clientId string) (*valueobject.Client, error) {
	hlog.Debug("clientRepository.GetClientById", "GetClientById", clientId)
	if client, ok := clients[clientId]; ok {
		return client, nil
	}
	client, err := r.keeperGateway.GetClientByPhoneNumberId(ctx, clientId)
	if err != nil {
		return nil, err
	}
	clients[clientId] = client
	return client, nil
}
