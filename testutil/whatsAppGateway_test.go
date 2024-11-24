package testutil

import (
	"context"
	"testing"

	"github.com/inter-hubly/linker/internal/gateway"
	"github.com/stretchr/testify/assert"
)

func TestWhatsAppGateway(t *testing.T) {
	app := gateway.NewWhatsApp()
	ctx := context.Background()
	t.Run("send message to whatsapp", func(t *testing.T) {
		messageTest := GetMessageTest()
		err := app.SendMessage(ctx, messageTest)
		assert.Nil(t, err)
	})
}
