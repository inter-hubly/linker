package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/domain/valueobject"
)

type Variables interface {
	GetVariablesByUserId(ctx context.Context, userId string, variables ...string) (map[string]interface{}, error)
}

type variablesRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

func NewVariables(ctx context.Context) *variablesRepository {

	var (
		once       sync.Once
		repository *variablesRepository
	)

	once.Do(func() {
		repository = &variablesRepository{
			connection: hmongo.GetConnection(ctx, "variables"),
		}
	})
	return repository
}

func (v *variablesRepository) GetVariablesByUserId(ctx context.Context, userId string, variables ...string) (map[string]interface{}, error) {
	field := valueobject.Pair[string, string]{Key: "userId", Value: userId}
	return v.connection.FindByFieldWithProjection(ctx, field, variables...)
}
