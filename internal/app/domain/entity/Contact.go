package entity

import (
	"github.com/google/uuid"
	"github.com/inter-hubly/pilot/domain/base"
)

type Contact struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
	base.Entity
}
