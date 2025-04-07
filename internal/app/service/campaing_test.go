package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/inter-hubly/linker/internal/app/repository/mocks"
	serviceMock "github.com/inter-hubly/linker/internal/app/service/mocks"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/domain/valueobject"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
)

type allMock struct {
	campaignRepository *mocks.MockCampaign
	contactRepository  *mocks.MockContact
	whatsAppService    *serviceMock.MockWhatsApp
}

func TestCampaign(t *testing.T) {
	ctx := testutils.SetLoggedUser(context.Background())
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	allMocks := allMock{
		campaignRepository: mocks.NewMockCampaign(ctrl),
		contactRepository:  mocks.NewMockContact(ctrl),
		whatsAppService:    serviceMock.NewMockWhatsApp(ctrl),
	}

	campaignServiceTest := campaignService{
		campaignRepository: allMocks.campaignRepository,
		contactRepository:  allMocks.contactRepository,
		whatsAppService:    allMocks.whatsAppService,
	}

	t.Run("should create a new campaign", func(t *testing.T) {
		// before
		allMocks.campaignRepository.EXPECT().GetCampaignById(gomock.Any(), gomock.Any()).Return(&entity.Campaign{
			Id:   "1234",
			Name: "test",
			Variables: []valueobject.Pair[string, string]{
				{"text", "name"},
			},
		}, nil)

		allMocks.contactRepository.EXPECT().GetContactsById(gomock.Any(), gomock.Any()).Return([]*entity.Contact{
			{
				Id:   "1234",
				Name: "test",
				Variables: []valueobject.Pair[string, string]{
					{"name", "testName"},
				},
			},
		}, nil)

		allMocks.whatsAppService.EXPECT().StartTemplate(gomock.Any(), gomock.Any()).Return(nil)

		err := campaignServiceTest.StartCampaign(ctx, "1234")
		assert.Nil(t, err)
	})
}
