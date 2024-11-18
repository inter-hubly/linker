package dto

type WhatsAppJSONReceived struct {
	ID       string              `json:"id"`
	Sender   WhatsAppPhoneIdDto  `json:"sender"`
	Receive  WhatsAppPhoneIdDto  `json:"receive"`
	Metadata WhatsAppMetadataDto `json:"metadata"`
}

type WhatsAppPhoneIdDto struct {
	PhoneNumberID      string `json:"phoneNumberId"`
	DisplayPhoneNumber string `json:"displayPhoneNumber"`
}

type WhatsAppMetadataDto struct {
	MessageID   string `json:"messageId"`
	RecipientID string `json:"recipientId"`
	Status      string `json:"status"`
	Body        string `json:"body"`
	Timestamp   int64  `json:"timestamp"`
	BodyLength  int64  `json:"bodyLength"`
}
