package express

import (
	"github.com/inter-hubly/linker/internal/controller"
	"github.com/inter-hubly/linker/internal/service"
	rabbitmq "github.com/inter-hubly/pilot/broker"
)

func Start() {
	rabbitmq.NewRabbitMQ("linker", "topic")

	service.NewWhatsApp()
	controller.NewWhatsApp()
}
