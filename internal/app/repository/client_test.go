package repository

import (
	"context"
	"testing"
	"time"

	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/testutils"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	host, close, err := testutils.Pgsql(ctx)
	if err != nil {
		panic(err)
	}
	defer close(ctx)

	if pgsql.GetConnection(ctx) == nil {
		pgsql.NewConnection(
			ctx,
			pgsql.WithUrl(host),
		)
	}

	repository := clientRepository{
		connection: pgsql.GetConnection(ctx),
	}

	for _, v := range []struct {
		testName string
		auxFunc  func()
	}{
		{
			testName: "Need to get a client",
			auxFunc: func() {
				insertedId := setUp(ctx, t, repository.connection)
				phoneNumberId := "5548"
				clientEntity, err := repository.GetClientByPhoneNumberId(ctx, phoneNumberId)
				assert.Nil(t, err)
				assert.Equal(t, insertedId, clientEntity.Id)
				assert.Equal(t, phoneNumberId, clientEntity.PhoneNumberId)
			},
		},
	} {
		t.Run(v.testName, func(t *testing.T) {
			v.auxFunc()
		})
	}
}

func setUp(ctx context.Context, t *testing.T, conn pgsql.SqlConn) string {
	exec, err := conn.Exec(`
			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
			CREATE TABLE IF NOT EXISTS client
			(
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				name            VARCHAR(255)        NOT NULL,
				email           VARCHAR(255) UNIQUE NOT NULL,
				app_id          VARCHAR(255)        NOT NULL,
				phone_number_id VARCHAR(255)    UNIQUE    NOT NULL,
				business_id     VARCHAR(255)        NOT NULL,
				access_token    VARCHAR(255)        NOT NULL,
				created_at      TIMESTAMP        NOT NULL,
				updated_at      TIMESTAMP        NOT NULL,
				removed          BOOL DEFAULT false
			);`)
	affected, err := exec.RowsAffected()
	if err != nil {
		panic(err)
	}
	assert.NotNil(t, affected)
	query := `INSERT INTO client (name, email, app_id, phone_number_id, business_id, access_token, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	row, err := conn.Query(query,
		"clientTest",
		"clientEmail",
		"12345",
		"5548",
		"23456",
		"12345asd",
		time.Now(),
		time.Now(),
	)
	assert.Nil(t, err)

	var returnedId string
	if err = row.Scan(&returnedId); err != nil {
		t.Fatal(err)
	}
	return returnedId
}
