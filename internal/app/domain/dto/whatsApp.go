package dto

type SenderMessageType string

const (
	TemplateMessageType SenderMessageType = "template"
	TextMessageType     SenderMessageType = "text"
)

type WhatsAppTextDto struct {
	PreviewUrl bool   `json:"preview_url"`
	Body       string `json:"body"`
}
