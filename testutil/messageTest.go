package testutil

import "github.com/inter-hubly/pilot/domain/entity"

func GetMessageTest() *entity.WhatsAppJSONReceived {
	return &entity.WhatsAppJSONReceived{
		Owner: entity.WhatsAppPhoneIdDto{
			PhoneNumberID:      "515719138282305",
			DisplayPhoneNumber: "15551817023",
		},
		SenderPhone: "+5548991784586",
		Metadata: entity.WhatsAppMetadataDto{
			Body: "Mensagem teste do aplicativo",
		},
	}
}
