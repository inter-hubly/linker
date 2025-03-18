package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCampaign(t *testing.T) {
	ctx := context.Background()
	campaignServiceTest := campaignService{}

	t.Run("should create a new campaign", func(t *testing.T) {
		err := campaignServiceTest.StartCampaign(ctx, "")
		assert.Nil(t, err)
	})
}
