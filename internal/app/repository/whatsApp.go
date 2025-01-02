package repository

import (
	"context"
	"sync"
	"time"

	"github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	"github.com/inter-hubly/linker/internal/app/domain/entity"
	"github.com/inter-hubly/pilot/database/elasticsearch"
)

type WhatsApp interface {
	PersistMessage(ctx context.Context, message *entity.Chat) (string, error)
	SetStatusMessageById(ctx context.Context, messageId string, status dto.MessageStatus, expirationTime int64) error
}

type whatsAppRepository struct {
	elastic elasticsearch.ElasticConn
	index   string
}

func NewWhatsApp() *whatsAppRepository {

	var (
		whatsAppOnce sync.Once
		whatsApp     *whatsAppRepository
	)

	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppRepository{
			elastic: elasticsearch.GetConnection(),
			index:   "whatsapp.ready",
		}
	})
	return whatsApp
}

func (w *whatsAppRepository) PersistMessage(ctx context.Context, message *entity.Chat) (string, error) {
	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	res, err := w.elastic.Create(ctx, w.index, message)
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func (w *whatsAppRepository) SetStatusMessageById(
	ctx context.Context,
	messageId string,
	status dto.MessageStatus,
	expirationTime int64,
) error {

	now := time.Now()
	query := map[string]interface{}{
		"script": map[string]interface{}{
			"source": `
            ctx._source.status.add(params.new_status); 
            ctx._source.expiration_time = params.new_expiration_time;
            ctx._source.updatedAt = params.updatedAt;
        `,
			"params": map[string]interface{}{
				"new_status": map[string]interface{}{
					"status":     status,
					"receivedAt": now.Unix(),
				},
				"new_expiration_time": expirationTime,
				"updatedAt":           now,
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
	_, err := w.elastic.Update(ctx, w.index, query)
	if err != nil {
		return err
	}
	return nil
}
