//go:build e2e

package testutil

func TestWhatsAppGateway(t *testing.T) {
	app := gateway.NewWhatsApp()
	ctx := context.Background()
	t.Run("send message to whatsapp", func(t *testing.T) {
		messageTest := GetMessageTest()

		messageDto := dto.GatewayWhatsAppMessageDto{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               messageTest.Sender.PhoneNumberId,
			Type:             "text",
			Text: &dto.WhatsAppTextDto{
				PreviewUrl: true,
				Body:       messageTest.Metadata.Body,
			},
		}

		_, err := app.SendMessage(ctx, messageTest.Owner.PhoneNumberId, &messageDto)
		assert.Nil(t, err)
	})
}
