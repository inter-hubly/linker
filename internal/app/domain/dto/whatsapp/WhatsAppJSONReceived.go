package dto

type OriginType string

const (
	UtilityOriginType OriginType = "utility"
)

type MessageType string

const (
	MessageTypeStatuses MessageType = "statuses"
	MessageTypeMessage  MessageType = "message"
)

type MessageStatus string

const (
	SentStatus      MessageStatus = "sent"
	DeliveredStatus MessageStatus = "delivered"
	ReadStatus      MessageStatus = "read"
	StartStatus     MessageStatus = "start"
	ErrorStatus     MessageStatus = "error"
)

type WhatsAppJSONReceived struct {
	Id          string              `json:"id,omitempty"`
	MessageType MessageType         `json:"messageType"`
	Owner       WhatsAppPhoneIdDto  `json:"owner,omitempty"`
	Sender      WhatsAppPhoneIdDto  `json:"sender,omitempty"`
	Status      MessageStatus       `json:"status,omitempty"`
	Metadata    WhatsAppMetadataDto `json:"metadata,omitempty"`
	Active      bool                `json:"active"`
}

type WhatsAppPhoneIdDto struct {
	PhoneNumberId string `json:"phoneNumberId"`
	PhoneNumber   string `json:"phoneNumber,omitempty"`
	ProfileName   string `json:"profileName,omitempty"`
}

type WhatsAppMetadataDto struct {
	ExpirationTimeStamp string     `json:"expirationTimeStamp,omitempty"`
	Timestamp           string     `json:"timestamp,omitempty"`
	ConversationId      string     `json:"conversationId,omitempty"`
	OriginType          OriginType `json:"originType,omitempty"`
	MessageId           string     `json:"messageId,omitempty"`
	Body                string     `json:"body,omitempty"`
	BodyLength          int        `json:"bodyLength,omitempty"`
}
