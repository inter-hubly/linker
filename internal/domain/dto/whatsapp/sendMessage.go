package dto

// SendTextDto is to send any message
type SendTextDto struct {
	SenderAndReceiver SenderAndReceiverDto `json:"senderAndReceiver"`
	Message           string               `json:"message"`
}
