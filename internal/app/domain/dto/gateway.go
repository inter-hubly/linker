package dto

type GatewayWhatsAppMessageDto struct {
	Text             *WhatsAppTextDto  `json:"text,omitempty"`
	Template         *TemplateBody     `json:"template,omitempty"`
	Type             SenderMessageType `json:"type"`
	MessagingProduct string            `json:"messaging_product"`
	To               string            `json:"to"`
	// RecipientType    string            `json:"recipient_type"`
}

// ResponseWhatsAppGateway is a Gateway response
type ResponseWhatsAppGateway struct {
	Contact          []Contact `json:"contacts"`
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

type TemplateBody struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}

type Language struct {
	Code string `json:"code"`
}

type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}
