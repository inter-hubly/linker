//go:build e2e

package controller

import (
	"context"
	"testing"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppController(t *testing.T) {
	for _, v := range []struct {
		testName string
	}{
		{
			testName: "Need pass in test because was marshaled",
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			ctx := context.Background()
			received, err := parseJsonReceived(ctx, []byte(jsonReceived))
			if err != nil {
				assert.NotNil(t, err)
			}

			assert.Equal(t, received.Id, "510006955530686")
			assert.Equal(t, received.Owner.PhoneNumberId, "515719138282305")
			assert.Equal(t, received.Status, dto.DeliveredStatus)
		})
	}
}

var jsonReceived = `{
  "id": "510006955530686",
  "messageType": "message",
  "owner": {
    "phoneNumberId": "515719138282305",
    "displayPhoneNumber": "15551817023"
  },
  "senderPhone": "554891784586",
  "status": "delivered",
  "metadata": {
    "timestamp": "1732229980",
    "messageId": "wamid.HBgMNTU0ODkxNzg0NTg2FQIAEhgWM0VCMEQxMTAyODM3RUM2RjM0OTlEOQA=",
    "body": "Olá! Sua mensagem teste foi recebida com sucesso!",
    "bodyLength": 50
  }
}`

// {
// "id":"67ec6ea0e58c687402ed07bf"
// }

// MENSAGEM DE RETORNO
// tenantId: 459185417288378
// {
// "id": "msg-123",
// "messageType": "text",
// "owner": {
// "phoneNumberId": "0987654321",
// "profileName": "Empresa XYZ"
// },
// "sender": {
// "phoneNumberId": "554891784586",
// "profileName": "Cliente João"
// },
// "status": "delivered",
// "metadata": {
// "expirationTimeStamp": "2025-04-06T00:00:00Z",
// "timestamp": "2025-04-05T23:00:00Z",
// "conversationId": "conv-abc-001",
// "originType": "user",
// "messageId": "wamid.HBgMNTU2NjY2NjY2NjY2FQIAEhggMTIzNDU2Nzg5MA==",
// "body": "Olá, tudo bem?",
// "bodyLength": 16
// },
// "active": true
// }
