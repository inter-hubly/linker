package domain

type WhatsApp struct {
	ID      string `json:"id"`
	Changes Change `json:"changes"`
}

type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}

type Value struct {
	Statuses         Status   `json:"statuses"`
	Metadata         Metadata `json:"metadata"`
	MessagingProduct string   `json:"messagingProduct"`
}

type Status struct {
	ID          string `json:"id"`
	RecipientID string `json:"recipientId"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
}
type Metadata struct {
	DisplayPhoneNumber string `json:"displayPhoneNumber"`
	PhoneNumberID      string `json:"phoneNumberID"`
}
