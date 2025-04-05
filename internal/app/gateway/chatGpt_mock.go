package gateway

import (
	"context"

	"github.com/inter-hubly/linker/internal/app/domain/dto"
)

type chatGptMock struct{}

func NewChatGptMock() Chatgpt {
	return &chatGptMock{}
}

func (c *chatGptMock) GetInformation(ctx context.Context, context *dto.IaContext, contexts []dto.IaContext) (string, error) {
	return "Chatgpt Mock", nil
}
