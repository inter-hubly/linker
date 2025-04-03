package repository

import (
	"testing"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/pilot/database/hredis"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/stretchr/testify/assert"
)

func TestIaContext(t *testing.T) {
	ctx := hctx.Tenant.New("1234556")
	hredis.NewConnection(ctx)
	iaContextTestRepository := NewIaContext(ctx)

	t.Run("save IaContextTestRepository", func(t *testing.T) {
		saveContext, err := iaContextTestRepository.SaveContext(ctx, "12345", &dto.IaContext{Content: "test"})
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
		err := iaContextTestRepository.StartContext(ctx, "123", "test")
		assert.Nil(t, err)
	})
}
