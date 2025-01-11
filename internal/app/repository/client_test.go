package repository

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	clients = make(map[string]*valueobject.Client)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	keeper := gateway.NewMockKeeper(ctrl)
	client := clientRepository{
		// connection: keeper,
	}
	returnedValue := &valueobject.Client{
		Id:            1,
		PhoneNumberId: "123",
		AccessToken:   "acess_token",
	}
	keeper.EXPECT().GetClient(
		gomock.Any(),
		gomock.Any(),
	).Return(
		returnedValue,
		nil,
	)
	ctx := context.Background()
	t.Run("Get Client in cache", func(t *testing.T) {
		firstValue, err := client.GetClientById(ctx, "1")
		assert.NoError(t, err)
		assert.NotEmpty(t, firstValue)
		assert.Equal(t, firstValue.Id, uint64(1))
		assert.Equal(t, firstValue.AccessToken, returnedValue.AccessToken)

		secondValue, err := client.GetClientById(ctx, "1")
		assert.NoError(t, err)
		assert.NotEmpty(t, secondValue)
		assert.Equal(t, secondValue.Id, uint64(1))
		assert.Equal(t, secondValue.AccessToken, returnedValue.AccessToken)
	})
}
