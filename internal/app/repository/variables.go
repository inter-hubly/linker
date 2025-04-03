//go:generate mockgen -source=variables.go -destination=mocks/variables_mock.go -package=mocks

package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Variables interface {
	GetVariablesByUserId(ctx context.Context, userId string, variables ...string) (map[string]interface{}, error)
}

type variablesRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

var (
	variablesOnce sync.Once
	variable      *variablesRepository
)

func NewVariables(ctx context.Context) *variablesRepository {
	variablesOnce.Do(func() {
		variable = &variablesRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "variables",
		}
	})
	return variable
}

func (v *variablesRepository) GetVariablesByUserId(ctx context.Context, userId string, variables ...string) (map[string]interface{}, error) {
	projection := bson.M{}
	for _, fld := range variables {
		projection[fld] = 1
	}

	opts := options.FindOne().SetProjection(projection)
	tenantId := hctx.Tenant.Get(ctx)
	result := make(map[string]interface{})
	if err := v.connection.GetCollection(ctx, v.collection).
		FindOne(ctx,
			bson.M{
				"userId":   userId,
				"tenantId": tenantId,
			},
			opts).
		Decode(&result); err != nil {
		hlog.Error(ctx, "variablesRepository.GetVariablesByUserId", err.Error())
		return nil, fmt.Errorf("failed to decode object: %w", err)
	}

	return result, nil
}
