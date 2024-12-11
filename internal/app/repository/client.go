package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/pilot/database/pgsql"
)

type Client interface {
	GetClientByPhoneId(ctx context.Context, numberId string) (string, error)
}

var (
	clientOnce sync.Once
	client     *clientRepository
)

type clientRepository struct {
	pgsqlConnection pgsql.SqlConn
}

func NewClient() *clientRepository {
	clientOnce.Do(func() {
		client = &clientRepository{
			// pgsqlConnection: pgsql.GetConnection(),
		}
	})
	return client
}

func (r *clientRepository) GetClientByPhoneId(ctx context.Context, numberId string) (string, error) {
	return map[string]string{
		"554896711701": "+5548996711701",
		"554891784586": "+5548991784586",
	}[numberId], nil
}
