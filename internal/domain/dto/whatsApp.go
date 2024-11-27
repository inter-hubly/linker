package dto

type SenderMessageType string

const (
	TemplateMessageType SenderMessageType = "template"
	TextMessageType     SenderMessageType = "text"
)

type GatewayWhatsAppMessageDto struct {
	Text             *WhatsAppTextDto  `json:"text,omitempty"`
	Template         *TemplateDto      `json:"template,omitempty"`
	MessagingProduct string            `json:"messaging_product"`
	RecipientType    string            `json:"recipient_type"`
	To               string            `json:"to"`
	Type             SenderMessageType `json:"type"`
}

type WhatsAppTextDto struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type StartTemplateDto struct {
	SenderAndReceiver SenderAndReceiverDto `json:"senderAndReceiver"`
	Name              string               `json:"name"`
	Language          string               `json:"language"`
}

type TemplateDto struct {
	Name        string      `json:"name"`
	LanguageDto LanguageDto `json:"language"`
}
type LanguageDto struct {
	Code string `json:"code"`
}

type Message struct {
	Id            string `json:"id"`
	MessageStatus string `json:"message_status"`
}

type SentTextDto struct {
	SenderAndReceiver SenderAndReceiverDto `json:"senderAndReceiver"`
	Message           string               `json:"message"`
}

type ResponseWhatsAppGateway struct {
	MessagingProduct string    `json:"messaging_product"`
	Contact          []Contact `json:"contact"`
	Messages         []Message `json:"messages"`
}

type Contact struct {
	Input string `json:"input"`
	WaId  string `json:"wa_id"`
}

type SenderAndReceiverDto struct {
	From string `json:"from"`
	To   string `json:"to"`
}
