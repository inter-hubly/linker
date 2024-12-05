package testutil

import (
	"encoding/json"

	dto2 "github.com/inter-hubly/linker/internal/app/domain/dto/whatsapp"
	entity2 "github.com/inter-hubly/linker/internal/app/domain/entity"
)

func GetMessageTest() *dto2.WhatsAppJSONReceived {
	return &dto2.WhatsAppJSONReceived{
		Owner: dto2.WhatsAppPhoneIdDto{
			PhoneNumberID:      "515719138282305",
			DisplayPhoneNumber: "15551817023",
		},
		SenderPhoneId: "+5548991784586",
		Metadata: dto2.WhatsAppMetadataDto{
			Body: "Mensagem teste do aplicativo",
		},
	}
}

func SentMessage() ([]byte, error) {
	res := dto2.SendTextDto{
		SenderAndReceiver: dto2.SenderAndReceiverDto{
			To:            "5548991784586",
			From:          "15551817023",
			OwnerNumberId: "515719138282305",
		},
		Message: "Messagem de teste",
	}
	marshal, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func GetStartTemplateMessage() dto2.StartTemplateDto {
	return dto2.StartTemplateDto{
		SenderAndReceiver: dto2.SenderAndReceiverDto{
			OwnerNumberId: "515719138282305",
			To:            "+5548991784586",
			From:          "15551817023",
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

func GetChatToSave(e *dto2.WhatsAppJSONReceived) *entity2.Chat {
	return &entity2.Chat{
		MessageId: e.Metadata.MessageId,
		OwnerId:   e.Owner.PhoneNumberID,
		// From:      e.Owner.DisplayPhoneNumber,
		// To:        e.SenderPhone,
		Message: "",
	}
}

func NewWhatsAppMessage() *dto2.WhatsAppJSONReceived {
	return &dto2.WhatsAppJSONReceived{
		Id:     "123456",
		Active: true,
		Owner: dto2.WhatsAppPhoneIdDto{
			PhoneNumberID:      "515719138282305",
			DisplayPhoneNumber: "15551817023",
		},
		Status:        dto2.DeliveredStatus,
		SenderPhoneId: "+5548991784586",
		Metadata: dto2.WhatsAppMetadataDto{
			ExpirationTimeStamp: "1732814400",
			MessageId:           "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSRjJCRDQ0NTIzNkNEMTY5Q0JGAA==",
			ConversationId:      "8496b7ae453f655b01ac1e5654713a1d",
			Timestamp:           "1731695647",
			OriginType:          dto2.UtilityOriginType,
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
