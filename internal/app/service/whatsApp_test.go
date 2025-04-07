package service

import (
	"testing"
)

func TestWhatsApp(t *testing.T) {
	// ctx := testutils.SetLoggedUser(context.Background())

	// whatsappServiceTest := NewWhatsApp(ctx)

	t.Run("Need to create context with flows", func(t *testing.T) {
		// err := whatsappServiceTest.ReceiveMessage(ctx, &dto.WhatsAppJSONReceived{
		//
		// })
		// assert.Nil(t, err)
	})
}

// {
// "to": "+5511999999999",
// "campaignId": "abc123",
// "templateInfo": {
// "id": "tpl-001",
// "name": "boas_vindas",
// "language": "pt_BR",
// "message": "Olá, seja bem-vindo!"
// },
// "hasIaInteraction": false,
// "parameters": [
// {
// "first": "nome",
// "second": "João"
// },
// {
// "first": "data",
// "second": "2025-04-05"
// }
// ]
// }
