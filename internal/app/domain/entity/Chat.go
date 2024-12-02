package entity

import (
	"github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
)

type ChatType string

const (
	ChatTemplate ChatType = "template"
	ChatText     ChatType = "text"
)

type Chat struct {
	Audit        []ChatMessageStatusTime `json:"status,omitempty"`
	Type         ChatType                `json:"type"`
	MessageId    string                  `json:"messageId,omitempty"`
	OwnerId      string                  `json:"ownerId"`
	OwnerPhone   string                  `json:"ownerPhone"`
	ToPhoneId    string                  `json:"toPhoneId"`
	ToPhone      string                  `json:"toPhone"`
	TemplateName string                  `json:"templateName,omitempty"`
	Message      string                  `json:"message,omitempty"`
}

type ChatMessageStatusTime struct {
	Status     dto.MessageStatus `json:"status"`
	ReceivedAt string            `json:"receivedAt"`
}
