package dto

type SendChatGpt struct {
	Model    string           `json:"model"`
	Messages []ChatGptMessage `json:"messages"`
}

type ChatGptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseChatGpt struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
