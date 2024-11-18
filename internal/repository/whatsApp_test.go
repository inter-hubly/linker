package repository

import (
	"context"
	"os"
	"testing"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
)

const environmentDatabase = false

func TestWhatsApp(t *testing.T) {
	ctx := context.Background()

	var host string
	var close func(ctx context.Context) error
	var err error

	if environmentDatabase {
		host, close, err = testutils.ElasticSearch(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if close != nil {
			defer close(ctx)
		}
	} else {
		os.Setenv("ENVIRONMENT", "test")
		server.MockStartEnv("../../")
		host = server.GetElasticSearch().Host
	}

	elasticsearch.NewConn(elasticsearch.WithUrl([]string{host}))
	repository := NewWhatsApp()
	for _, v := range []struct {
		testName string
	}{
		{
			testName: "have to persist",
		},
	} {
		t.Run(v.testName, func(t *testing.T) {

			err := repository.PersistMessage(ctx, NewWhatsAppMessage())
			assert.Nil(t, err)
		})
	}
}

func NewWhatsAppMessage() *dto.WhatsAppJSONReceived {
	return &dto.WhatsAppJSONReceived{
		ID: "123456",
		Sender: dto.WhatsAppPhoneIdDto{
			PhoneNumberID:      "515719138282305",
			DisplayPhoneNumber: "15551817023",
		},
		Receive: dto.WhatsAppPhoneIdDto{
			PhoneNumberID:      "510006955530686",
			DisplayPhoneNumber: "554891784586",
		},
		Metadata: dto.WhatsAppMetadataDto{
			MessageID:   "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSOTQyQjZBNEEwRjg3N0VGRURDAA==",
			RecipientID: "554891784586",
			Status:      "read",
			Timestamp:   1731695647,
		},
	}
}
