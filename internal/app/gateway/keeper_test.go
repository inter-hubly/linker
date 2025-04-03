//go:build e2e

package gateway

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeeper(t *testing.T) {
	kpGateway := NewKeeper()
	ctx := context.Background()
	t.Run("Need get client", func(t *testing.T) {
		client, err := kpGateway.GetClientByPhoneNumberId(ctx, "559153210606318")
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}
