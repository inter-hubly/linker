package testutil

import (
	"context"
	"testing"

	dto2 "github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	"github.com/inter-hubly/linker/internal/app/gateway"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppGateway(t *testing.T) {
	app := gateway.NewWhatsApp()
	ctx := context.Background()
	t.Run("send message to whatsapp", func(t *testing.T) {
		messageTest := GetMessageTest()

		messageDto := dto2.GatewayWhatsAppMessageDto{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               messageTest.SenderPhoneId,
			Type:             "text",
			Text: &dto2.WhatsAppTextDto{
				PreviewUrl: true,
				Body:       messageTest.Metadata.Body,
			},
		}

		_, err := app.SendMessage(ctx, messageTest.Owner.PhoneNumberID, &messageDto)
		assert.Nil(t, err)
	})
}
