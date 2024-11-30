package testutil

import (
	"testing"

	rabbitmq "github.com/inter-hubly/pilot/broker"
)

func TestEndToEnd(t *testing.T) {
	rabbitmq.NewRabbitMQ("linker", "topic")
	message, err := StartMessage()
	if err != nil {
		t.Fatal(err)
	}

	rabbitmq.GetConnection().Publish("whatsapp.start", message)

}
