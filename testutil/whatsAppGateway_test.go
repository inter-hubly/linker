package testutil

import (
	"context"
	"testing"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/linker/internal/gateway"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppGateway(t *testing.T) {
	app := gateway.NewWhatsApp()
	ctx := context.Background()
	t.Run("send message to whatsapp", func(t *testing.T) {
		messageTest := GetMessageTest()

		messageDto := dto.GatewayWhatsAppMessageDto{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               messageTest.SenderPhone,
			Type:             "text",
			Text: &dto.WhatsAppTextDto{
				PreviewUrl: true,
				Body:       messageTest.Metadata.Body,
			},
		}

		_, err := app.SendMessage(ctx, messageTest.Owner.PhoneNumberID, &messageDto)
		assert.Nil(t, err)
	})
}
