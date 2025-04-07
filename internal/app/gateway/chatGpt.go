package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
	"github.com/inter-hubly/pilot/domain/entity"
)

type Chatgpt interface {
	GetInformation(context.Context, *entity.Flow, []entity.Flow) (string, error)
}

type chatgpt struct {
	basicUrl string
	token    string
	model    string
	context  map[string]string
}

func NewChatgpt(ctx context.Context) Chatgpt {
	return &chatgpt{
		basicUrl: "https://api.openai.com/v1/chat/completions",
		token:    "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",
		model:    "gpt-4o-mini",
		context:  make(map[string]string),
	}
}

func (c *chatgpt) GetInformation(ctx context.Context, iaMessage *entity.Flow, iaContext []entity.Flow) (string, error) {

	iaContextMessages := make([]dto.ChatGptMessage, 0, len(iaContext))
	for _, v := range iaContext {
		gptMessage := dto.ChatGptMessage{
			Role:    v.Role,
			Content: v.Message,
		}
		iaContextMessages = append(iaContextMessages, gptMessage)
	}

	iaContextMessages = append(iaContextMessages, dto.ChatGptMessage{
		Role:    iaMessage.Role,
		Content: iaMessage.Message,
	})

	gptBody := dto.SendChatGpt{
		Model:    c.model,
		Messages: iaContextMessages,
	}
	marshal, err := json.Marshal(gptBody)
	if err != nil {
		return "", err
	}
	body := bytes.NewReader(marshal)

	// TODO setar no padrao da api
	request, err := http.NewRequest(http.MethodPost, c.basicUrl, body)
	if err != nil {
		return "", err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New(response.Status)
	}

	var myResp dto.ResponseChatGpt
	if err = json.Unmarshal(bodyBytes, &myResp); err != nil {
		return "", fmt.Errorf("error decoding JSON into entity.Company: %w. Raw response: %s", err, string(bodyBytes))
	}

	return myResp.Choices[0].Message.Content, nil
}
