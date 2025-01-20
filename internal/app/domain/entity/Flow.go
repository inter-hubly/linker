package entity

import "github.com/inter-hubly/pilot/domain/base"

type FlowType string

const (
	Messaging         FlowType = "messaging"
	ProcessActivation FlowType = "processActivation"
)

type Flow struct {
	ID        string   `bson:"id"`
	Name      string   `bson:"name"`
	Text      string   `bson:"text"`
	NextSteps string   `bson:"steps"`
	FlowType  FlowType `bson:"flow_type"`
	base.Entity
}
