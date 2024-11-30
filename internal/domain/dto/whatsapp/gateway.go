package dto

type GatewayWhatsAppMessageDto struct {
	Text             *WhatsAppTextDto  `json:"text,omitempty"`
	Template         *TemplateDto      `json:"template,omitempty"`
	MessagingProduct string            `json:"messaging_product"`
	RecipientType    string            `json:"recipient_type"`
	To               string            `json:"to"`
	Type             SenderMessageType `json:"type"`
}

// ResponseWhatsAppGateway is a Gateway response
type ResponseWhatsAppGateway struct {
	MessagingProduct string    `json:"messaging_product"`
	Contact          []Contact `json:"contact"`
	Messages         []Message `json:"messages"`
}

type Contact struct {
	Input string `json:"input"`
	WaId  string `json:"wa_id"`
}

type Message struct {
	Id            string `json:"id"`
	MessageStatus string `json:"message_status"`
}
