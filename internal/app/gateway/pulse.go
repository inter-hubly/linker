package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/pilot/server"
)

type Pulse interface {
	HandleMessage(ctx context.Context, ownerId string, message *dto.PulseDto) error
}

var (
	pulseOnce sync.Once
	pulse     *pulseGateway
)

type pulseGateway struct {
	url string
}

func NewPulse() *pulseGateway {
	pulseOnce.Do(func() {
		pulse = &pulseGateway{
			url: server.GetGatewayHost().PulseHost,
		}
	})
	return pulse
}

func (p *pulseGateway) HandleMessage(ctx context.Context, ownerId string, message *dto.PulseDto) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/receive?user=%s", p.url, ownerId), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
