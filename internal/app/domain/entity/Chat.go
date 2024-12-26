package entity

import (
	"time"

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
	ProfileName  string                  `json:"profileName,omitempty"`
	OwnerId      string                  `json:"ownerId"`
	ToPhone      string                  `json:"toPhone,omitempty"`
	TemplateName string                  `json:"templateName,omitempty"`
	Message      string                  `json:"message,omitempty"`
	IsOwner      bool                    `json:"isOwner"`
	CreatedAt    time.Time               `json:"createdAt,omitempty"`
	UpdatedAt    time.Time               `json:"updatedAt,omitempty"`
}

type ChatMessageStatusTime struct {
	Status     dto.MessageStatus `json:"status"`
	ReceivedAt int64             `json:"receivedAt"`
}
