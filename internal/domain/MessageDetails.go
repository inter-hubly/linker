package domain

import "github.com/inter-hubly/pilot/domain/valueobject"

type MessageDetails struct {
	Template  string         `json:"template"`
	Id        valueobject.Id `json:"id"`
	MessageId valueobject.Id `json:"messageId"`
	ClientId  valueobject.Id `json:"clientId"`
	UserId    valueobject.Id `json:"userId"`
}
