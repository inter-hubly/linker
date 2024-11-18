package gateway

import (
	"context"
	"fmt"
	"github.com/inter-hubly/pilot/hlog"
	"sync"

	"github.com/inter-hubly/linker/internal/domain/dto"
)

type WhatsApp interface {
	GetAccessToken(ctx context.Context)
	SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
	ReceiveMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error
}

var (
	whatsAppOnce sync.Once
	whatsApp     *whatsAppGateway
)

type whatsAppGateway struct {
	url string
}

func NewWhatsApp() *whatsAppGateway {
	whatsAppOnce.Do(func() {
		whatsApp = &whatsAppGateway{
			url: "https://graph.facebook.com/v21.0/",
		}
	})
	return whatsApp
}

func (w *whatsAppGateway) SendMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {
	hlog.Debug("whatsAppGateway.SendMessage", fmt.Sprintf("Send Message %s", message))
	//url := fmt.Sprintf(w.url, "416481324892567/messages")
	//accessToken := "EAAPdvbIZBdiIBOZCZADTUqZCelwE9E6WIP4KnPq5zhsqmiPmxjcES9mVc1Xrqz8b5sEGHuCSsGM5U16NQfwnxc8GwZBAVoJrWMZCcY93ZAJgDZAV7viYiAdZAvxN1c77tpY0jsBWINUPKGESOddlf0VZAuUZAZANS3JlUitYjpoZBG8bMDwaG3gxmGcZBPsGOPGYThuBWDDZCmB2sQZC56xBCjPcnicTyRkp"
	//
	//payload := []byte("{ \"messaging_product\": \"whatsapp\", \"to\": \"5548991784586\", \"type\": \"template\", \"template\": { \"name\": \"hello_world\", \"language\": { \"code\": \"en_US\" } } }")
	//
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	//if err != nil {
	//	log.Fatalf("Erro ao criar requisição: %v", err)
	//}
	//
	//// Configurar os headers
	//req.Header.Set("Authorization", "Bearer "+accessToken)
	//req.Header.Set("Content-Type", "application/json")
	//
	//// Enviar a requisição
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Fatalf("Erro ao enviar requisição: %v", err)
	//}
	//defer resp.Body.Close()
	//
	//// Validar resposta
	//if resp.StatusCode != http.StatusOK {
	//	log.Printf("Erro na resposta: StatusCode = %d", resp.StatusCode)
	//} else {
	//	log.Println("Mensagem enviada com sucesso!")
	//}
	//
	//fmt.Println("Resposta da API:", resp.Status)
	//all, err := io.ReadAll(resp.Body)
	//fmt.Println("body da API: ", string(all))
	//return nil
	return nil
}

func (w *whatsAppGateway) GetAccessToken(ctx context.Context) {

}

func (w *whatsAppGateway) ReceiveMessage(ctx context.Context, message *dto.WhatsAppJSONReceived) error {
	return nil
}
