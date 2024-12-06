package dto

type GatewayWhatsAppMessageDto struct {
	Text             *WhatsAppTextDto  `json:"text,omitempty"`
	Template         *TemplateDto      `json:"template,omitempty"`
	Type             SenderMessageType `json:"type"`
	MessagingProduct string            `json:"messaging_product"`
	To               string            `json:"to"`
	RecipientType    string            `json:"recipient_type"`
}

// ResponseWhatsAppGateway is a Gateway response
type ResponseWhatsAppGateway struct {
	Contact          []Contact `json:"contact"`
	Messages         []Message `json:"messages"`
	MessagingProduct string    `json:"messaging_product"`
}

type Contact struct {
	Input string `json:"input"`
	WaId  string `json:"wa_id"`
}

type Message struct {
	Id            string `json:"id"`
	MessageStatus string `json:"message_status"`
}
