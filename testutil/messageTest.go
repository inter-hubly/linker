package testutil

import (
	"encoding/json"

	dto "github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	"github.com/inter-hubly/linker/internal/app/domain/entity"
)

func GetMessageTest() *dto.WhatsAppJSONReceived {
	return &dto.WhatsAppJSONReceived{
		Owner: dto.WhatsAppPhoneIdDto{
			PhoneNumber: "515719138282305",
		},
		Sender: dto.WhatsAppPhoneIdDto{
			PhoneNumber: "554891784586",
		},
		Metadata: dto.WhatsAppMetadataDto{
			Body: "Mensagem teste do aplicativo",
		},
	}
}

func SentMessage() ([]byte, error) {
	res := dto.SendTextDto{
		SenderAndReceiver: dto.SenderAndReceiverDto{
			To: "5548991784586",
		},
		Message: "Messagem de teste",
	}
	marshal, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func GetStartTemplateMessage() dto.StartTemplateDto {
	return dto.StartTemplateDto{
		SenderAndReceiver: dto.SenderAndReceiverDto{
			To: "+5548991784586",
		},
		Name:     "hello_world",
		Language: "en_US",
	}
}

func StartMessage() ([]byte, error) {
	res := GetStartTemplateMessage()
	marshal, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func GetChatToSave(e *dto.WhatsAppJSONReceived) *entity.Chat {
	return &entity.Chat{
		MessageId: e.Metadata.MessageId,
		OwnerId:   e.Owner.PhoneNumber,
		Message:   "",
	}
}

func NewWhatsAppMessage() *dto.WhatsAppJSONReceived {
	return &dto.WhatsAppJSONReceived{
		Id:     "123456",
		Active: true,
		Owner: dto.WhatsAppPhoneIdDto{
			PhoneNumber: "515719138282305",
		},
		Status: dto.DeliveredStatus,
		Sender: dto.WhatsAppPhoneIdDto{
			PhoneNumber: "554891784586",
		},
		Metadata: dto.WhatsAppMetadataDto{
			ExpirationTimeStamp: "1732814400",
			MessageId:           "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSRjJCRDQ0NTIzNkNEMTY5Q0JGAA==",
			ConversationId:      "8496b7ae453f655b01ac1e5654713a1d",
			Timestamp:           "1731695647",
			OriginType:          dto.UtilityOriginType,
		},
	}
}

// {
// "senderAndReceiver" : {
// "from" : "15551817023",
// "to" : "+5548991784586"
// },
// "name" : "hello_world",
// "language" : "en_US"
// }
