package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/stretchr/testify/assert"
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
			var entityWhats dto.WhatsAppJSONReceived
			err := json.Unmarshal([]byte(jsonReceived), &entityWhats)
			assert.Nil(t, err)

			err = service.SendMessage(ctx, &entityWhats)
			assert.Nil(t, err)
		})
	}
}

var jsonReceived = `
{
	"id": "510006955530686",
	"sender": {
			"phoneNumberId": "515719138282305",
			"displayPhoneNumber": "15551817023"
	},
	"receive": {
			"phoneNumberId": "510006955530686",
			"displayPhoneNumber": "554891784586"
	},
	"metadata": {
			"messageId": "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSOTQyQjZBNEEwRjg3N0VGRURDAA==",
			"recipientId": "554891784586",
			"status": "read",
			"body": "",
			"timestamp": 1731695647
	}
}`
