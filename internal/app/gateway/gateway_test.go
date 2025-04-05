//go:build e2e

package gateway

import (
	"encoding/json"
	"testing"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {

	var res dto.ResponseWhatsAppGateway
	err := json.Unmarshal(getResponse(), &res)
	assert.NoError(t, err)
	assert.Equal(t, res.Contact[0].Input, "5548991784586")
	assert.Equal(t, res.Messages[0].Id, "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSQjI5QzI5MjI1REU5QzZDQjI2AA==")
}

func getResponse() []byte {
	return []byte(` 
		{
		  "messaging_product" : "whatsapp",
		  "contacts" : [ {
			"input" : "5548991784586",
			"wa_id" : "554891784586"
		  } ],
		  "messages" : [ {
			"id" : "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSQjI5QzI5MjI1REU5QzZDQjI2AA==",
			"message_status" : "accepted"
		  } ]
		}`)
}
