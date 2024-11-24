package dto

type GatewayWhatsAppMessageDto struct {
	Text             WhatsAppTextDto `json:"text"`
	MessagingProduct string          `json:"messaging_product"`
	RecipientType    string          `json:"recipient_type"`
	To               string          `json:"to"`
	Type             string          `json:"type"`
}

type WhatsAppTextDto struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}
