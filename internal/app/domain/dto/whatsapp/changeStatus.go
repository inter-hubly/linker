package dto

// ChangeStatusDto is about status
type ChangeStatusDto struct {
	Status              MessageStatus `json:"status"`
	MessageId           string        `json:"messageId"`
	ExpirationTimeStamp int64         `json:"expirationTimeStamp,omitempty"`
}
