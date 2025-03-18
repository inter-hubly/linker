package repository

import (
	"context"
	"testing"

	"github.com/inter-hubly/pilot/database/hredis"
)

func TestIaContext(t *testing.T) {
	ctx := context.Background()
	
	hredis.NewConnection(ctx)
	iaContextTestRepository := NewIaContext(ctx)
	
	
	t.Run("test IaContextTestRepository", func(t *testing.T) {
		iaContextTestRepository.
	})
	
}
