package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	"github.com/inter-hubly/pilot/hlog"
)

type WhatsApp interface {
	GetAccessToken(ctx context.Context)
	SendMessage(ctx context.Context, phoneNumberId string, message *dto.GatewayWhatsAppMessageDto) (*dto.ResponseWhatsAppGateway, error)
	ReadyMessage(ctx context.Context, phoneNumberId, messageId string) error
}

var (
	whatsAppOnce  sync.Once
	whatsApp      *whatsAppGateway
	graphAPIToken = "EAAP8VwxXKXkBOZCAGZBXwMg96AWefXAxCGH4vNMQEYzvQ4KYH4B7NTE39c8A5nqFNTJPwOzZBh2ZAmrDoejTwZBhX4Q1jL7GgJnqVxutKbwUobFAspao8JecTZCeQtExCD9my2SPlznZCBt4KnrgIcJRR8UvX67ZCxZAr4xwvP57Eae8ad9zxW9xLjnJGaGDuYtRxTDZA7nPMsHpKsvxgUN5g8AjsQ"
	// graphAPIToken = "1"
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
) (*dto.ResponseWhatsAppGateway, error) {
	hlog.Debug("whatsAppGateway.SendMessage", fmt.Sprintf("Send Message %v", messageDto))
	res, err := w.makeRequest(http.MethodPost, fmt.Sprintf("%s%s/messages", w.url, phoneNumberId), messageDto)

	// TODO testar
	if err != nil {
		hlog.Error("whatsAppGateway.SendMessage", fmt.Sprintf("Failed to send message %v", messageDto))
		return nil, err
	}
	return res, nil
}

func (w *whatsAppGateway) GetAccessToken(ctx context.Context) {

}

func (w *whatsAppGateway) ReadyMessage(ctx context.Context, phoneNumberId, messageId string) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageId,
	}
	if _, err := w.makeRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s/messages", w.url, phoneNumberId),
		data,
	); err != nil {
		return err
	}
	return nil
}

func (w *whatsAppGateway) makeRequest(method, url string, data any) (*dto.ResponseWhatsAppGateway, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+graphAPIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP request returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var res dto.ResponseWhatsAppGateway

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &res, err
}
