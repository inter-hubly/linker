package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/inter-hubly/linker/internal/app/domain/dto"
)

type ChatType string

const (
	ChatTemplate ChatType = "template"
	ChatText     ChatType = "text"
)

type Chat struct {
	Audit       []ChatMessageStatusTime `json:"status,omitempty"`
	Type        ChatType                `json:"type"`
	MessageId   string                  `json:"messageId,omitempty"`
	ProfileName string                  `json:"profileName,omitempty"`
	OwnerId     string                  `json:"ownerId"`
	ToPhoneId   string                  `json:"toPhoneId,omitempty"`
	CampaignId  uuid.UUID               `json:"campaignId,omitempty"`
	Message     string                  `json:"message,omitempty"`
	IsOwner     bool                    `json:"isOwner"`
	CreatedAt   time.Time               `json:"createdAt,omitempty"`
	UpdatedAt   time.Time               `json:"updatedAt,omitempty"`
}

type ChatMessageStatusTime struct {
	Status     dto.MessageStatus `json:"status"`
	ReceivedAt int64             `json:"receivedAt"`
}
