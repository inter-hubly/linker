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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetCampaign(t *testing.T) {
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

	repository := campaignRepository{
		connection: hmongo.GetConnection(ctx),
		collection: "campaigns",
	}

	for _, v := range []struct {
		testName string
		auxFunc  func()
	}{
		{
			testName: "Need to get a campaign",
			auxFunc: func() {
				campaignEntity := entity.Campaign{
					Name: "Campaign Test",
				}
				loggedUser := hctx.LoggedUser.Get(ctx)
				campaignEntity.Entity = base.NewBaseEntity(ctx, &loggedUser)
				result, err := repository.connection.GetCollection(ctx, repository.collection).InsertOne(ctx, campaignEntity)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotNil(t, result.InsertedID)

				id, err := repository.GetCampaignById(ctx, result.InsertedID.(primitive.ObjectID).Hex())
				assert.Nil(t, err)
				assert.Equal(t, campaignEntity.Name, id.Name)
			},
		},
		{
			testName: "Cant get a campaign with other tenant",
			auxFunc: func() {
				campaignEntity := entity.Campaign{
					Name: "Campaign Test",
				}
				campaignEntity.Entity = base.NewBaseEntity(ctx, &hctx.Logged{
					Tenant: "wrongTenant",
				})

				result, err := repository.connection.GetCollection(ctx, repository.collection).InsertOne(ctx, campaignEntity)
				if err != nil {
					t.Fatal(err)
				}
				assert.NotNil(t, result.InsertedID)

				_, err = repository.GetCampaignById(ctx, result.InsertedID.(primitive.ObjectID).Hex())
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, "mongo: no documents in result")
				assert.NotEqual(t, hctx.Tenant.Get(ctx), "")
				var foundValue entity.Campaign
				newErr := repository.connection.
					GetCollection(ctx, repository.collection).FindOne(ctx,
					bson.M{
						"_id": result.InsertedID.(primitive.ObjectID),
					}).Decode(&foundValue)
				assert.Nil(t, newErr)
				assert.Equal(t, foundValue.Name, campaignEntity.Name)
				assert.Equal(t, foundValue.TenantId, campaignEntity.TenantId)
			},
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			v.auxFunc()
		})
	}
}
