package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/inter-hubly/linker/internal/domain/dto/whatsapp"
	entity2 "github.com/inter-hubly/linker/internal/domain/entity"
	"github.com/inter-hubly/linker/testutil"
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
		auxFunc  func() (string, error)
	}{
		{
			testName: "have to persist",
			auxFunc: func() (string, error) {
				chatToSave := testutil.GetChatToSave(testutil.NewWhatsAppMessage())

				chatToSave.Audit = append(chatToSave.Audit, entity2.ChatMessageStatusTime{
					Status:     dto.DeliveredStatus,
					ReceivedAt: fmt.Sprintf("%d", time.Now().Unix()),
				})
				return repository.PersistMessage(ctx, chatToSave)
			},
		},
		{
			testName: "Need change status",
			auxFunc: func() (string, error) {
				chatToSave := testutil.GetChatToSave(testutil.NewWhatsAppMessage())

				chatToSave.Audit = append(chatToSave.Audit, entity2.ChatMessageStatusTime{
					Status:     dto.DeliveredStatus,
					ReceivedAt: fmt.Sprintf("%d", time.Now().Unix()),
				})

				chatId, err := repository.PersistMessage(ctx, chatToSave)
				assert.Nil(t, err)
				assert.NotEmpty(t, chatId)

				err = repository.SetStatusMessageById(ctx, chatToSave.MessageId, dto.SentStatus)
				assert.Nil(t, err)

				resp, err := repository.elastic.FindById(ctx, "whatsapp.ready", chatId)
				assert.Nil(t, err)
				statusResp := resp.Source["status"].([]interface{})
				assert.Equal(t, len(statusResp), 2)
				return chatId, nil
			},
		},
	} {
		t.Run(v.testName, func(t *testing.T) {

			id, err := v.auxFunc()
			assert.Nil(t, err)
			assert.NotEmpty(t, id)
		})
	}
}
