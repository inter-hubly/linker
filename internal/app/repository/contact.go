//go:generate mockgen -source=contact.go -destination=mocks/contact_mock.go -package=mocks

package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/inter-hubly/pilot/database/hmongo"
	"github.com/inter-hubly/pilot/domain/entity"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact interface {
	GetContactsById(ctx context.Context, contactsId ...string) ([]*entity.Contact, error)
}

type contactRepository struct {
	connection hmongo.NoSqlConn
	collection string
}

var (
	contactOnce sync.Once
	contact     *contactRepository
)

func NewContact(ctx context.Context) *contactRepository {

	contactOnce.Do(func() {
		contact = &contactRepository{
			connection: hmongo.GetConnection(ctx),
			collection: "contact",
		}
	})
	return contact
}

func (r *contactRepository) GetContactsById(ctx context.Context, contactsId ...string) ([]*entity.Contact, error) {
	hlog.Debug(ctx, "contactRepository.GetContactsById", fmt.Sprint("contactsId", contactsId))
	tenantId := hctx.Tenant.Get(ctx)

	ids := make([]primitive.ObjectID, 0, len(contactsId))
	for _, id := range contactsId {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			hlog.Error(ctx, "contactRepository.GetContactsById.err", err.Error())
			return nil, err
		}
		ids = append(ids, objId)
	}

	filter := bson.M{
		"tenantId": tenantId,
		"_id": bson.M{
			"$in": ids,
		},
	}

	cur, err := r.connection.GetCollection(ctx, r.collection).Find(ctx, filter)
	if err != nil {
		hlog.Error(ctx, "contactRepository.GetContactsById", fmt.Sprint("error finding contacts: ", err.Error()))
		return nil, err
	}

	var contacts []*entity.Contact
	for cur.Next(ctx) {
		var ct entity.Contact
		if err = cur.Decode(&ct); err != nil {
			hlog.Error(ctx, "contactRepository.GetContactsById", fmt.Sprint("error decoding ", err.Error()))
			continue
		}
		contacts = append(contacts, &ct)
	}

	if err = cur.Err(); err != nil {
		hlog.Error(ctx, "contactRepository.GetContactsById", fmt.Sprint("error iterating cursor: ", err.Error()))
		return nil, err
	}

	return contacts, nil
}
