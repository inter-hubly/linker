package repository

import (
	"context"
	"sync"
	"time"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/domain/entity"
)

type WhatsApp interface {
	PersistMessage(ctx context.Context, message *entity.Chat) (string, error)
	SetStatusMessageById(ctx context.Context, messageId string, status entity.MessageStatus) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppRepository
	elasticIndex = "whatsapp.ready"
)

type whatsAppRepository struct {
	elastic elasticsearch.ElasticConn
}

func NewWhatsApp() *whatsAppRepository {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppRepository{
			elastic: elasticsearch.GetConnection(),
		}
	})
	return whatsApp
}

func (w *whatsAppRepository) PersistMessage(ctx context.Context, message *entity.Chat) (string, error) {
	res, err := w.elastic.Create(ctx, elasticIndex, message)
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func (w *whatsAppRepository) SetStatusMessageById(
	ctx context.Context,
	messageId string,
	status entity.MessageStatus,
) error {
	query := map[string]interface{}{
		"script": map[string]interface{}{
			"source": "ctx._source.status.add(params.new_status);",
			"params": map[string]interface{}{
				"new_status": map[string]interface{}{
					"status":     status,
					"receivedAt": time.Now().Unix(),
				},
			},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"messageId": messageId,
						},
					},
				},
			},
		},
	}
	_, err := w.elastic.Update(ctx, elasticIndex, query)
	if err != nil {
		return err
	}
	return nil
}
