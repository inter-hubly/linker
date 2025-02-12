package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

var clients map[string]*entity.Client

type Client interface {
	GetClientById(ctx context.Context, clientId string) (*entity.Client, error)
}

type clientRepository struct {
	connection pgsql.SqlConn
}

var (
	clientOnce sync.Once
	client     *clientRepository
)

func NewClient() *clientRepository {
	clients = make(map[string]*entity.Client)

	clientOnce.Do(func() {
		client = &clientRepository{
			connection: pgsql.GetConnection(),
		}
	})
	return client
}

func (r *clientRepository) GetClientById(ctx context.Context, clientId string) (*entity.Client, error) {
	hlog.Debug(ctx, "clientRepository.GetClientById", fmt.Sprint("GetClientById", clientId))
	if client, ok := clients[clientId]; ok {
		return client, nil
	}
	query := `SELECT c.id, c.name, c.email, c.app_id, c.phone_number_id, c.business_id, c.access_token 
          FROM client c 
          WHERE c.phone_number_id = $1`

	queryExec, err := r.connection.Query(query, clientId)
	if err != nil {
		hlog.Error(ctx, "clientRepository.GetClientById", fmt.Sprintf("error find clientId %s : %s", clientId, err))
		return nil, err
	}
	var clientDb entity.Client
	if err = queryExec.Scan(
		&clientDb.Id,
		&clientDb.Name,
		&clientDb.Email,
		&clientDb.AppId,
		&clientDb.PhoneNumberId,
		&clientDb.BusinessId,
		&clientDb.AccessToken,
	); err != nil {
		hlog.Error(ctx, "clientRepository.GetClientById", fmt.Sprintf("error scan clientId %s : %s", clientId, err))
		return nil, err
	}
	return &clientDb, nil
}
