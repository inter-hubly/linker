package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/linker/internal/app/repository"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/hrest"
)

type WhatsApp interface {
	SendMessage(ctx context.Context, phoneNumberId string, message *dto.GatewayWhatsAppMessageDto) (*dto.ResponseWhatsAppGateway, error)
	ReadyMessage(ctx context.Context, phoneNumberId, messageId string) error
}

type whatsAppGateway struct {
	url              string
	clientRepository repository.Client
}

func NewWhatsApp() *whatsAppGateway {
	var (
		whatsAppOnce sync.Once
		whatsApp     *whatsAppGateway
	)

	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppGateway{
			url:              "https://graph.facebook.com/v21.0/",
			clientRepository: repository.NewClient(),
		}
	})
	return whatsApp
}

func (w *whatsAppGateway) SendMessage(
	ctx context.Context,
	phoneNumberId string,
	messageDto *dto.GatewayWhatsAppMessageDto,
) (*dto.ResponseWhatsAppGateway, error) {
	hlog.Debug(ctx, "whatsAppGateway.SendMessage", fmt.Sprintf("Send Message %v", messageDto))

	client, err := w.clientRepository.GetClientById(ctx, phoneNumberId)
	if err != nil {
		return nil, err
	}

	request := hrest.NewRequest(fmt.Sprintf("%s%s/messages", w.url, phoneNumberId),
		hrest.WithHeader([]hrest.Pair[string, string]{
			{"Content-Type", "application/json"},
			{"Authorization", "Bearer " + client.AccessToken},
		}),
		hrest.WithBody(messageDto),
	)

	if err = request.CreateRequest(ctx, http.MethodPost); err != nil {
		hlog.Error(ctx, "whatsAppGateway.SendMessage", fmt.Sprintf("Failed to send message %v", messageDto))
		return nil, err
	}

	var resp dto.ResponseWhatsAppGateway
	if err = request.GetBody(ctx, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (w *whatsAppGateway) ReadyMessage(ctx context.Context, phoneNumberId, messageId string) error {
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageId,
	}
	if _, err := w.makeRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s%s/messages", w.url, phoneNumberId),
		data,
	); err != nil {
		return err
	}
	return nil
}

func (w *whatsAppGateway) makeRequest(ctx context.Context, method, url string, data any) (*dto.ResponseWhatsAppGateway, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+" w.token")
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
