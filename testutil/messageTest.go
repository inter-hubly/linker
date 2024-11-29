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

func GetChatToSave(e *entity.WhatsAppJSONReceived) *entity.Chat {
	return &entity.Chat{
		MessageId: e.Metadata.MessageId,
		OwnerId:   e.Owner.PhoneNumberID,
		From:      e.Owner.DisplayPhoneNumber,
		To:        e.SenderPhone,
		Message:   "",
	}
}

func NewWhatsAppMessage() *entity.WhatsAppJSONReceived {
	return &entity.WhatsAppJSONReceived{
		Id:     "123456",
		Active: true,
		Owner: entity.WhatsAppPhoneIdDto{
			PhoneNumberID:      "515719138282305",
			DisplayPhoneNumber: "15551817023",
		},
		Status:      entity.DeliveredStatus,
		SenderPhone: "+5548991784586",
		Metadata: entity.WhatsAppMetadataDto{
			ExpirationTimeStamp: "1732814400",
			MessageId:           "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSRjJCRDQ0NTIzNkNEMTY5Q0JGAA==",
			ConversationId:      "8496b7ae453f655b01ac1e5654713a1d",
			Timestamp:           "1731695647",
			OriginType:          entity.UtilityOriginType,
		},
	}
}
