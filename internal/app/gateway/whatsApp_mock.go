package gateway

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/hrest"
	"github.com/inter-hubly/pilot/server"
)

type whatsAppGatewayMock struct {
	url string
}

func NewWhatsAppMock() *whatsAppGatewayMock {
	var (
		whatsAppOnce sync.Once
		whatsApp     *whatsAppGatewayMock
	)

	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppGatewayMock{
			url: server.GetGatewayHost().WebhookHost,
		}
	})
	return whatsApp
}

func (w whatsAppGatewayMock) SendMessage(ctx context.Context, phoneNumberId string, messageDto *dto.GatewayWhatsAppMessageDto) (*dto.ResponseWhatsAppGateway, error) {
	hlog.Debug(ctx, "whatsAppGatewayMock.SendMessage", fmt.Sprintf("Send Mock Message %v", messageDto))

	smsBody := "test"
	if messageDto != nil && messageDto.Text != nil {
		smsBody = messageDto.Text.Body
	}
	entryMessage := newWhatsAppEntry(smsBody, messageDto.To, phoneNumberId, messageDto.To, phoneNumberId)
	request := hrest.NewRequest(fmt.Sprintf("%s/webhook", w.url),
		hrest.WithBody(entryMessage),
	)

	if err := request.CreateRequest(ctx, http.MethodPost); err != nil {
		hlog.Error(ctx, "whatsAppGatewayMock.SendMessage", fmt.Sprintf("Failed to send message %v", messageDto))
		return nil, err
	}

	return &dto.ResponseWhatsAppGateway{
		Messages: []dto.Message{
			{Id: "test1"},
		},
	}, nil
}

func (w whatsAppGatewayMock) ReadyMessage(ctx context.Context, phoneNumberId, messageId string) error {
	// TODO implement me
	panic("implement me")
}

type whatsAppEntry struct {
	Object string  `json:"object"`
	Entry  []entry `json:"entry"`
}

type entry struct {
	Id      string   `json:"id,omitempty"`
	Changes []change `json:"changes"`
}

type change struct {
	Value value  `json:"value"`
	Field string `json:"field"`
}

type value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         metadata  `json:"metadata,omitempty"`
	Messages         []message `json:"messages,omitempty"`
	Contacts         []contact `json:"contacts,omitempty"`
}

type metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type message struct {
	From      string         `json:"from"`
	Id        string         `json:"id"`
	Timestamp string         `json:"timestamp"`
	Text      messageContent `json:"text"`
	Type      string         `json:"type"`
}

type messageContent struct {
	Body string `json:"body"`
}

type contact struct {
	Profile profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type profile struct {
	Name string `json:"name"`
}

func newWhatsAppEntry(msgText, waId, from, to, tenantId string) whatsAppEntry {
	return whatsAppEntry{
		Object: "whatsapp_business_account",
		Entry: []entry{
			{
				Id: "525894273941090",
				Changes: []change{
					{
						Value: value{
							MessagingProduct: "whatsapp",
							Metadata: metadata{
								DisplayPhoneNumber: to,
								PhoneNumberID:      tenantId,
							},
							Messages: []message{
								{
									From:      from,
									Id:        "wamid.HBgMNTU0ODkxNzg0NTg2FQIAEhgWM0VCMDYzQ0FCMjIzMkJCNDJBMEEzQgA=",
									Timestamp: "1736213333",
									Text:      messageContent{Body: msgText},
									Type:      "text",
								},
							},
							Contacts: []contact{
								{
									Profile: profile{
										Name: "Test User",
									},
									WaID: waId,
								},
							},
						},
						Field: "messages",
					},
				},
			},
		},
	}
}
