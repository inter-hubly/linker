package lentity

import "github.com/inter-hubly/pilot/domain/base"

type Variables struct {
	Id     string `bson:"_id"`
	UserId string `bson:"userId"`
	base.Entity
}
