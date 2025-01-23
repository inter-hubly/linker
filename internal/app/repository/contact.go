package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/inter-hubly/linker/internal/app/domain/entity"
	"github.com/inter-hubly/pilot/database/pgsql"
	"github.com/inter-hubly/pilot/hlog"
)

type Contact interface {
	GetContactsById(ctx context.Context, contactsId ...string) ([]*entity.Contact, error)
}

type contactRepository struct {
	connection pgsql.SqlConn
}

func NewContact() *contactRepository {

	var (
		clientOnce sync.Once
		client     *contactRepository
	)

	clientOnce.Do(func() {
		client = &contactRepository{
			connection: pgsql.GetConnection(),
		}
	})
	return client
}

func (r *contactRepository) GetContactsById(ctx context.Context, contactsId ...string) ([]*entity.Contact, error) {
	hlog.Debug(ctx, "contactRepository.GetContactsById", fmt.Sprint("contactsId", contactsId))

	placeholders := make([]string, len(contactsId))
	args := make([]interface{}, len(contactsId))
	for i, id := range contactsId {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	sqlQuery := fmt.Sprintf("SELECT c.id, c.name, c.phone FROM contact c WHERE c.id IN (%s)", strings.Join(placeholders, ", "))
	query, err := r.connection.QueryRows(sqlQuery, args...)
	hlog.Info(ctx, "contactRepository.GetContactsById", fmt.Sprintf("sqlQuery %s %s", sqlQuery, args))
	if err != nil {
		hlog.Error(ctx, "contactRepository.GetContactsById", fmt.Sprintf("error in query %s", err.Error()))
		return nil, err
	}
	defer query.Close()

	var contacts []*entity.Contact
	for query.Next() {
		var contact entity.Contact
		if err = query.Scan(&contact.Id, &contact.Name, &contact.Phone); err != nil {
			hlog.Error(ctx, "contactRepository.GetContactsById", fmt.Sprintf("error scanning row: %s", err.Error()))
			return nil, err
		}
		contacts = append(contacts, &contact)
	}

	return contacts, nil
}
