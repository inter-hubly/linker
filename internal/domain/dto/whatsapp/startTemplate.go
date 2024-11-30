package dto

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
