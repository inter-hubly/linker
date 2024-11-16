package service

import (
	"context"
	"testing"
)

func TestWhatsApp(t *testing.T) {
	service := NewWhatsApp()

	for _, v := range []struct {
		testName string
		auxFunc  func()
	}{
		{
			testName: "need to send message",
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			ctx := context.Background()
			service.SendMessage(ctx, []byte(jsonData))
		})
	}
}

var jsonData = `
{
	"id": "510006955530686",
	"changes": [
		{
			"field": "messages",
			"value": {
				"messaging_product": "whatsapp",
				"metadata": {
					"display_phone_number": "15551817023",
					"phone_number_id": "515719138282305"
				},
				"statuses": [
					{
						"id": "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSOTQyQjZBNEEwRjg3N0VGRURDAA==",
						"recipient_id": "554891784586",
						"status": "read",
						"timestamp": 1731695647
					}
				]
			}
		}
	]
}
`
