package cache

import (
	"testing"

	"github.com/inter-hubly/pilot/database/hredis"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
)

func TestIaContext(t *testing.T) {
	ctx := hctx.Tenant.New("1234556")

	if hredis.GetConnection(ctx) == nil {
		redis, close, err := testutils.Redis(ctx)
		defer close(ctx)
		assert.NoError(t, err)

		hredis.NewConnection(ctx, hredis.WithAddr(redis))
	}

	iaContextTestRepository := flowContextRepository{
		connection: hredis.GetConnection(ctx),
	}

	t.Run("save IaContextTestRepository", func(t *testing.T) {
		saveContext, err := iaContextTestRepository.SaveContext(ctx, "12345", &entity.Flow{Message: "test", IsIaInteraction: true})
		assert.Nil(t, err)
		assert.NotNil(t, saveContext)
	})

	t.Run("get all IaContextTestRepository", func(t *testing.T) {
		saveContext, err := iaContextTestRepository.GetContext(ctx, "12345")
		assert.Nil(t, err)
		assert.NotEmpty(t, saveContext)
	})

	t.Run("get all IaContextTestRepository was error", func(t *testing.T) {
		getAllMessages, err := iaContextTestRepository.GetContext(ctx, "12")
		assert.NotNil(t, err)
		assert.Empty(t, getAllMessages)
	})

	t.Run("start IaContextTestRepository", func(t *testing.T) {
		err := iaContextTestRepository.StartContext(ctx, "123", &entity.Flow{Message: "test"})
		assert.Nil(t, err)
	})
}
