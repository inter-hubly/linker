package express

import (
	"github.com/inter-hubly/linker/internal/app/controller"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
)

const ExchangeBroker = "linker"

func Start() {

	rabbitmq.NewRabbitMQ(ExchangeBroker, "topic")

	if err := rabbitmq.GetConnection().
		QueueBind(
			rabbitmq.NewQueueBinding("whatsapp.start", "whatsapp.start", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.statuses", "whatsapp.statuses", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.message", "whatsapp.message", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.sent", "whatsapp.sent", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.send", "whatsapp.send", ExchangeBroker),
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
