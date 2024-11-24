package testutil

import "github.com/inter-hubly/pilot/domain/dto"

func GetMessageTest() *dto.WhatsAppJSONReceived {
	return &dto.WhatsAppJSONReceived{
		Owner: dto.WhatsAppPhoneIdDto{
			PhoneNumberID:      "515719138282305",
			DisplayPhoneNumber: "15551817023",
		},
		SenderPhone: "+5548991784586",
		Metadata: dto.WhatsAppMetadataDto{
			Body: "Mensagem teste do aplicativo",
		},
	}
}
