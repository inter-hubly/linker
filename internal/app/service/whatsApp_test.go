package service

import (
	"context"
	"testing"

	"github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	"github.com/inter-hubly/pilot/server"
	"github.com/stretchr/testify/assert"
)

func TestWhatsApp(t *testing.T) {
	server.NewMockEnvironment(
		server.MockEnvironment{
			WhatsAppToken: "EAAP8VwxXKXkBO3WJuZAX2ZBsW2mVzBAHqG7gb38x2fIgd3ydReFv6VVNRZBGnZBsmhOt8CF7R9msJgPSprnjS7dM7g6XQdfvo50dguTVp6NxZAzC6JX5SvDqqBv3CpHX39O1uEAkkCdBrbwE4NKbCBHKVNOC7dtc2OSFQWSnIx2HUbvCTFFbg3lXxFt3j4PQZA",
		},
	)

	service := NewWhatsApp()

	for _, v := range []struct {
		testName string
		auxFunc  func()
	}{
		{
			testName: "need to send message",
		},
	} {
		t.Run(v.testName, func(t *testing.T) {

			ctx := context.Background()
			err := service.SendMessage(ctx, &dto.SendTextDto{
				Message: "Test message",
				SenderAndReceiver: dto.SenderAndReceiverDto{
					OwnerId: "559153210606318",
					To:      "554888356622",
				},
			})
			assert.Nil(t, err)
		})
	}
}
