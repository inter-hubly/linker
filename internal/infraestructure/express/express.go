package express

import (
	"context"

	"github.com/inter-hubly/linker/internal/app/controller"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/server"
)

const ExchangeBroker = "linker"

func Start(ctx context.Context) {

	rabbitmq.NewRabbitMQ(ctx, ExchangeBroker, "topic", rabbitmq.WithURL(server.GetAmpqConfig().Host))

	if err := rabbitmq.GetConnection().
		QueueBind(
			ctx,
			rabbitmq.NewQueueBinding("whatsapp.start", "whatsapp.start", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.statuses", "whatsapp.statuses", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.message", "whatsapp.message", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.sent", "whatsapp.sent", ExchangeBroker),
			rabbitmq.NewQueueBinding("whatsapp.send", "whatsapp.send", ExchangeBroker),
			rabbitmq.NewQueueBinding("campaign.init", "campaign.init", ExchangeBroker),
		); err != nil {
		panic(err)
	}

	pgsql.NewConnection(
		pgsql.WithUrl(server.GetPgsqlConfig().Host),
	)

	hmongo.NewConnection(
		ctx,
		hmongo.WithDatabase(server.GetMongoConfig().Database),
		hmongo.WithUrl(server.GetMongoConfig().Host),
	)

	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{server.GetElasticSearch().Host}),
		elasticsearch.WithUsernameAndPassword(
			server.GetElasticSearch().Username,
			server.GetElasticSearch().Password,
		),
	)

	controller.NewWhatsApp(ctx)
	controller.NewFlow(ctx)
}
