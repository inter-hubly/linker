package repository

import (
	"context"
	"testing"

	"github.com/inter-hubly/pilot/database/elasticsearch"
	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/testutils"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	closeMongo := mongoSetup(ctx)
	closePgsql := pgsqlSetup(ctx)
	closeElastic := elasticSetup(ctx)
	defer closeMongo(ctx)
	defer closePgsql(ctx)
	defer closeElastic(ctx)
	m.Run()
}

func mongoSetup(ctx context.Context) func(context.Context) error {
	mongoHost, close, err := testutils.Mongo(ctx)
	if err != nil {
		panic(err)
	}

	hmongo.NewConnection(
		ctx,
		hmongo.WithDatabase("test"),
		hmongo.WithUrl(mongoHost),
	)
	return close
}

func pgsqlSetup(ctx context.Context) func(context.Context) error {
	pgsqlHost, close, err := testutils.Pgsql(ctx)
	if err != nil {
		panic(err)
	}
	pgsql.NewConnection(
		ctx,
		pgsql.WithUrl(pgsqlHost),
	)
	return close
}

func elasticSetup(ctx context.Context) func(context.Context) error {
	elasticHost, close, err := testutils.ElasticSearch(ctx)
	if err != nil {
		panic(err)
	}
	elasticsearch.NewConn(
		elasticsearch.WithUrl([]string{elasticHost}),
	)
	return close
}
