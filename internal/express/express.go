package express

import (
	"github.com/inter-hubly/linker/internal/controller"
	rabbitmq "github.com/inter-hubly/pilot/broker"
)

func Start() {
	rabbitmq.NewRabbitMQ("linker", "topic")

	controller.NewWhatsApp()
}
