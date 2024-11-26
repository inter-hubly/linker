package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	GetAccessToken(ctx context.Context)
	SendMessage(ctx context.Context, phoneNumberId string, message *dto.GatewayWhatsAppMessageDto) error
	ReceiveMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error
	ReadyMessage(ctx context.Context, phoneNumberId, messageId string) error
}

var (
	whatsAppOnce  sync.Once
	whatsApp      *whatsAppGateway
	graphAPIToken = "EAAP8VwxXKXkBO9gQBYWkLswZCFJk8ZCiFAjTOTBXJyU3HnjZA6mZAnHchT1PEtvmByg2ID4PrnoZBtVUjGv8o7aaqEIxpInrIOMMdDctGMKIS4BBhhZC9w02lztHNwSMEBYroJJgqRhj8eaSP6NVcCpuKQGcKxpyU6xWqmvpeQ6PqZBHiGq9NyxHuD67fVpTUsrdWk8WbeouVUI7lRDlM6emWFQ7wZDZD"
)

type whatsAppGateway struct {
	url string
}

func NewWhatsApp() *whatsAppGateway {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppGateway{
			url: "https://graph.facebook.com/v21.0/",
		}
	})
	return whatsApp
}

func (w *whatsAppGateway) SendMessage(
	_ context.Context,
	phoneNumberId string,
	messageDto *dto.GatewayWhatsAppMessageDto,
) error {
	hlog.Debug("whatsAppGateway.SendMessage", fmt.Sprintf("Send Message %v", messageDto))
	if err := w.makeRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s/messages", w.url, phoneNumberId),
		messageDto,
	); err != nil {
		hlog.Error("whatsAppGateway.SendMessage", fmt.Sprintf("Failed to send message %v", messageDto))
		return err
	}
	return nil
}

func (w *whatsAppGateway) GetAccessToken(ctx context.Context) {

}

func (w *whatsAppGateway) ReceiveMessage(ctx context.Context, message *entity.WhatsAppJSONReceived) error {
	return nil
}

func (w *whatsAppGateway) ReadyMessage(ctx context.Context, phoneNumberId, messageId string) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageId,
	}
	return w.makeRequest(http.MethodPost, fmt.Sprintf("%s%s/messages", w.url, phoneNumberId), data)
}

func (w *whatsAppGateway) makeRequest(method, url string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+graphAPIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP request returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
