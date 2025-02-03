package mediator

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/inter-hubly/linker/testutil"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppMediator(t *testing.T) {
	ctrl := gomock.NewController(t)
	whatsGateway := gateway.NewMockWhatsApp(ctrl)
	// whatsRepository := repository.NewMockWhatsApp(ctrl)
	mediator := whatsAppMediator{
		whatsAppGateway:    whatsGateway,
		whatsAppRepository: nil,
		// whatsAppRepository: whatsRepository,
	}

	for _, v := range []struct {
		testeName string
		auxFunc   func()
	}{
		{
			testeName: "Need to start template",
			auxFunc: func() {
				ctx := context.Background()
				whatsGateway.EXPECT().SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(&dto.ResponseWhatsAppGateway{
					Contact: []dto.Contact{
						{
							Input: "input",
							WaId:  "123456",
						},
					},
					Messages: []dto.Message{
						{
							MessageStatus: "accepted",
						},
					},
				}, nil)

				// whatsRepository.EXPECT().
				// 	PersistMessage(gomock.Any(), gomock.Any()).
				// 	Return("", nil)

				message := testutil.GetStartTemplateMessage()
				err := mediator.StartTemplate(ctx, &message)
				assert.NoError(t, err)
			},
		},
	} {
		t.Run(v.testeName, func(t *testing.T) {
			v.auxFunc()
		})
	}
}
