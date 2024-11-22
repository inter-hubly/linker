package express

import (
	"github.com/inter-hubly/linker/internal/controller"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
)

const ExchangeBroker = "linker"

func Start() {

	rabbitmq.NewRabbitMQ(ExchangeBroker, "topic")

	if err := rabbitmq.GetConnection().
		QueueBind(
			rabbitmq.NewQueueBinding("whatsapp.read", "whatsapp.read", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.sent", "whatsapp.sent", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.delivered", "whatsapp.delivered", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.received", "whatsapp.received", ExchangeBroker),
		); err != nil {
		panic(err)
	}

	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{server.GetElasticSearch().Host}),
		elasticsearch.WithUsernameAndPassword(
			server.GetElasticSearch().Username,
			server.GetElasticSearch().Password,
		),
	)

	controller.NewWhatsApp()
}
