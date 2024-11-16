package gateway

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func receive() {
	url := "https://graph.facebook.com/v21.0/416481324892567/messages"
	accessToken := "EAAPdvbIZBdiIBOZCZADTUqZCelwE9E6WIP4KnPq5zhsqmiPmxjcES9mVc1Xrqz8b5sEGHuCSsGM5U16NQfwnxc8GwZBAVoJrWMZCcY93ZAJgDZAV7viYiAdZAvxN1c77tpY0jsBWINUPKGESOddlf0VZAuUZAZANS3JlUitYjpoZBG8bMDwaG3gxmGcZBPsGOPGYThuBWDDZCmB2sQZC56xBCjPcnicTyRkp"

	payload := []byte("{ \"messaging_product\": \"whatsapp\", \"to\": \"5548991784586\", \"type\": \"template\", \"template\": { \"name\": \"hello_world\", \"language\": { \"code\": \"en_US\" } } }")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	// Configurar os headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Enviar a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	// Validar resposta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro na resposta: StatusCode = %d", resp.StatusCode)
	} else {
		log.Println("Mensagem enviada com sucesso!")
	}

	fmt.Println("Resposta da API:", resp.Status)
	all, err := io.ReadAll(resp.Body)
	fmt.Println("body da API: ", string(all))
}
