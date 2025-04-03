package repository

import (
	"context"
	"testing"

	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetCampaign(t *testing.T) {
	ctx := context.Background()
	host, close, err := testutils.Mongo(ctx)
	if err != nil {
		panic(err)
	}
	hmongo.NewConnection(
		ctx,
		hmongo.WithDatabase("test"),
		hmongo.WithUrl(host),
	)
	defer close(ctx)

	repository := campaignRepository{
		connection: hmongo.GetConnection(ctx),
		collection: "campaigns",
	}

	for _, v := range []struct {
		testName string
	}{
		{
			testName: "Need to get a campaign",
		},
	} {
		t.Run(v.testName, func(t *testing.T) {

			campaignEntity := entity.Campaign{
				Name: "Campaign Test",
			}
			result, err := repository.connection.GetCollection(ctx, repository.collection).InsertOne(ctx, campaignEntity)
			if err != nil {
				t.Fatal(err)
			}
			assert.NotNil(t, result.InsertedID)

			id, err := repository.GetCampaignById(ctx, result.InsertedID.(primitive.ObjectID).Hex())
			assert.Nil(t, err)
			assert.Equal(t, campaignEntity.Name, id.Name)
		})
	}
}
