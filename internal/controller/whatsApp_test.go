package controller

import (
	"context"
	"testing"

	"github.com/inter-hubly/pilot/domain/entity"
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
			assert.Equal(t, received.Owner.PhoneNumberID, "515719138282305")
			assert.Equal(t, received.Owner.DisplayPhoneNumber, "15551817023")
			assert.Equal(t, received.Status, entity.DeliveredStatus)
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
