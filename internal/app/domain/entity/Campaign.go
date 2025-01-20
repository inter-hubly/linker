package entity

import "github.com/inter-hubly/pilot/domain/base"

type Campaign struct {
	Id               string   `bson:"id"`
	Name             string   `bson:"name"`
	TemplateName     string   `bson:"templateName"`
	TemplateLanguage string   `bson:"templateLanguage"`
	Phones           []string `bson:"phones"`
	Flows            []string `bson:"flows"`
	ParametersLength uint8    `bson:"parametersLength"`
	base.Entity
}
