package express

import (
	"github.com/inter-hubly/linker/internal/controller"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/server"
)

func Start() {
	rabbitmq.NewRabbitMQ("linker", "topic")

	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{server.GetElasticSearch().Host}),
		elasticsearch.WithUsernameAndPassword(
			server.GetElasticSearch().Username,
			server.GetElasticSearch().Password,
		),
	)

	controller.NewWhatsApp()
}
