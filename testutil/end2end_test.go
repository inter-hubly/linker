package testutil

import (
	"github.com/inter-hubly/linker/internal/express"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/server"
	"sync"
	"testing"
	"time"
)

const withTestContainer = false

func TestEnd2End(t *testing.T) {
	var wg sync.WaitGroup
	if withTestContainer {

	} else {
		wg.Add(1)
		go prepareTestEnvironment()
	}

	time.Sleep(5 * time.Second)
	rabbitmq.GetConnection().Publish("whatsapp.sent", []byte(jsonReceived))
	time.Sleep(5 * time.Second)
	wg.Done()
}

func prepareTestEnvironment() {
	server.MockStartEnv("../")
	express.Start()
}

var jsonReceived = `{
  "id": "510006955530686",
  "messageType": "message",
  "owner": {
    "phoneNumberId": "515719138282305",
    "displayPhoneNumber": "15551817023"
  },
  "senderPhone": "554891784586",
  "status": "delivered",
  "metadata": {
    "timestamp": "1732229980",
    "messageId": "wamid.HBgMNTU0ODkxNzg0NTg2FQIAEhgWM0VCMEQxMTAyODM3RUM2RjM0OTlEOQA=",
    "body": "Olá! Sua mensagem teste foi recebida com sucesso!",
    "bodyLength": 50
  }
}`
