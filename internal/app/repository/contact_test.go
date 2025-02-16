package repository

import (
	"context"
	"testing"

	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestFindManyContact(t *testing.T) {
	ctx := context.Background()
	pgsql.NewConnection(pgsql.WithUrl("postgres://postgres:frajolinha202@localhost:5432/hubly?sslmode=disable"))
	ct := NewContact(ctx)

	contacts, err := ct.GetContactsById(ctx, "1", "2")
	assert.Nil(t, err)
	assert.Equal(t, "48991784586", contacts[0].Phone)
}
