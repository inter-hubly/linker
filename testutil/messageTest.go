package testutil

import (
	"encoding/json"

	"github.com/inter-hubly/linker/internal/domain/dto"
	"github.com/inter-hubly/pilot/domain/entity"
)

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

func SentMessage() ([]byte, error) {
	res := dto.SentTextDto{
		SenderAndReceiver: dto.SenderAndReceiverDto{
			To:   "5548991784586",
			From: "15551817023",
		},
		Message: "Messagem de teste",
	}
	marshal, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func StartMessage() ([]byte, error) {
	res := dto.StartTemplateDto{
		SenderAndReceiver: dto.SenderAndReceiverDto{
			To:   "+5548991784586",
			From: "15551817023",
		},
		Name:     "hello_world",
		Language: "en_US",
	}
	marshal, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
