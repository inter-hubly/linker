package gateway

import (
	"context"

	"github.com/inter-hubly/pilot/domain/entity"
)

type chatGptMock struct{}

func NewChatGptMock() Chatgpt {
	return &chatGptMock{}
}

func (c *chatGptMock) GetInformation(ctx context.Context, context *entity.Flow, contexts []entity.Flow) (string, error) {
	return "Chatgpt Mock", nil
}
