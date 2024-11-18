package repository

import (
	"context"
	"sync"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/pilot/database/elasticsearch"
)

type WhatsApp interface {
	PersistMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
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

func (w *whatsAppRepository) PersistMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {
	_, err := w.elastic.Create(ctx, elasticIndex, message)
	if err != nil {
		return err
	}
	return nil
}
