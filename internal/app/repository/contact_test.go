package repository

import (
	"context"
	"testing"

	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFindManyContact(t *testing.T) {
	ctx := testutils.SetLoggedUser(context.Background())
	host, close, err := testutils.Mongo(ctx)
	if err != nil {
		panic(err)
	}
	defer close(ctx)

	if hmongo.GetConnection(ctx) == nil {
		hmongo.NewConnection(
			ctx,
			hmongo.WithDatabase("test"),
			hmongo.WithUrl(host),
		)
	}

	repository := contactRepository{
		connection: hmongo.GetConnection(ctx),
		collection: "contact",
	}

	for _, v := range []struct {
		testName string
		auxFunc  func()
	}{
		{
			testName: "Need to get a contact",
			auxFunc: func() {
				contactEntity := entity.Contact{
					Name: "test",
				}
				loggedUser := hctx.LoggedUser.Get(ctx)
				contactEntity.Entity = base.NewBaseEntity(ctx, &loggedUser)

				one, err := repository.connection.GetCollection(ctx, repository.collection).InsertOne(ctx, contactEntity)
				assert.Nil(t, err)
				assert.NotNil(t, one)

				contactDatabase, err := repository.GetContactsById(ctx, one.InsertedID.(primitive.ObjectID).Hex())
				assert.Nil(t, err)
				assert.Equal(t, contactEntity.Name, contactDatabase[0].Name)
			},
		},
		{
			testName: "cant find a contact because the tenant is wrong",
			auxFunc: func() {
				contactEntity := entity.Contact{
					Name: "test wrong",
				}

				contactEntity.Entity = base.NewBaseEntity(ctx, &hctx.Logged{
					Tenant: "wrongTenant",
				})

				one, err := repository.connection.GetCollection(ctx, repository.collection).InsertOne(ctx, contactEntity)
				assert.Nil(t, err)
				assert.NotNil(t, one)

				response, err := repository.GetContactsById(ctx, one.InsertedID.(primitive.ObjectID).Hex())
				assert.Nil(t, err)
				assert.Empty(t, response)
				assert.NotEqual(t, hctx.Tenant.Get(ctx), "")
			},
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			v.auxFunc()
		})
	}
}
