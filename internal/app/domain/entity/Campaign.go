package entity

import (
	"github.com/google/uuid"
	"github.com/inter-hubly/pilot/domain/base"
	"github.com/inter-hubly/pilot/domain/valueobject"
)

type Campaign struct {
	Id         uuid.UUID                          `bson:"id"`
	Name       string                             `bson:"name"`
	Template   base.TemplateInfo                  `bson:"template"`
	ContactsId []string                           `bson:"contactsId"`
	Parameters []valueobject.Pair[string, string] `bson:"parameters"`
	base.Entity
}
