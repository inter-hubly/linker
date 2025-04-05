package lentity

type Flow struct {
	Id         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	HasAiCheck bool   `json:"hasAiCheck,omitempty" bson:"hasAiCheck,omitempty"`
	Message    string `json:"message,omitempty" bson:"message,omitempty"`
}
