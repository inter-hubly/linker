package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/inter-hubly/linker/internal/domain/entity"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/domain/entity"
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
			message := NewWhatsAppMessage()

			id, err := repository.PersistMessage(ctx, entity.NormalizeWhatsAppMessage(message))
			assert.Nil(t, err)

			chatMessageTime := entity.ChatMessageTime{
				CreatedInDatabase: time.Now(),
				ReceivedAt:        message.Metadata.Timestamp,
			}

			repository.SetStatusMessageById(ctx, message.Metadata.MessageId, entity.ReadStatus, chatMessageTime)

			chat, err := repository.elastic.FindById(ctx, "whatsapp.ready", id)

			chatMessage := chat.Source["read"].(map[string]interface{})

			assert.Nil(t, err)
			assert.NotNil(t, chat)
			assert.Equal(t, chat.ID, id)
			parsedTime, err := time.Parse(time.RFC3339Nano, chatMessage["CreatedInDatabase"].(string))
			if err != nil {
				return
			}
			assert.True(t, chatMessageTime.CreatedInDatabase.Equal(parsedTime))
			assert.Equal(t, chatMessageTime.ReceivedAt, chatMessage["ReceivedAt"])

		})
	}
}

func NewWhatsAppMessage() *entity.WhatsAppJSONReceived {
	return &entity.WhatsAppJSONReceived{
		Id: "123456",
		Owner: entity.WhatsAppPhoneIdDto{
			PhoneNumberID:      "515719138282305",
			DisplayPhoneNumber: "15551817023",
		},
		Status: entity.DeliveredStatus,
		Metadata: entity.WhatsAppMetadataDto{
			MessageId:      "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSOTQyQjZBNEEwRjg3N0VGRURDAA==",
			ConversationId: "read",
			Timestamp:      "1731695647",
			Body:           "ok",
		},
	}
}
