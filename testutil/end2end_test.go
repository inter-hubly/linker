package testutil

import (
	"sync"
	"testing"
	"time"

	"github.com/inter-hubly/linker/internal/express"
	rabbitmq "github.com/inter-hubly/pilot/broker"
	"github.com/inter-hubly/pilot/server"
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
	message, err := StartMessage()
	if err != nil {
		t.Fatal(err)
	}
	rabbitmq.GetConnection().Publish("whatsapp.start", message)
	wg.Wait()
	// time.Sleep(5 * time.Second)
	// wg.Done()
}

func prepareTestEnvironment() {
	server.MockStartEnv("../")
	express.Start()
}

var jsonReceived = map[string][]string{
	"whatsapp.delivered": {`{
  "id": "510006955530686",
  "messageType": "statuses",
  "owner": {
    "phoneNumberId": "515719138282305",
    "displayPhoneNumber": "15551817023"
  },
  "senderPhone": "554891784586",
  "status": "delivered",
  "metadata": {
    "timestamp": "1732229947",
    "conversationId": "7afc5c0fbaf6f5777c7be1b40b13c243",
    "messageId": "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSRDk2MDc1RjU4OUYwRTgxNzdDAA=="
  }
}`,
	},
	"whatsapp.read": {
		`{
  "id": "510006955530686",
  "messageType": "statuses",
  "owner": {
    "phoneNumberId": "515719138282305",
    "displayPhoneNumber": "15551817023"
  },
  "senderPhone": "554891784586",
  "status": "read",
  "metadata": {
    "timestamp": "1732229971",
    "messageId": "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSRDk2MDc1RjU4OUYwRTgxNzdDAA=="
  }
}`,
	},
	"whatsapp.sent": {
		`{
  "id": "510006955530686",
  "messageType": "statuses",
  "owner": {
    "phoneNumberId": "515719138282305",
    "displayPhoneNumber": "15551817023"
  },
  "senderPhone": "+5548991784586",
  "status": "sent",
  "metadata": {
    "timestamp": "1732229947",
    "conversationId": "7afc5c0fbaf6f5777c7be1b40b13c243",
    "messageId": "wamid.HBgMNTU0ODkxNzg0NTg2FQIAERgSRDk2MDc1RjU4OUYwRTgxNzdDAA==",
	"body":"isso é uma mensagem teste"
  }
}`},
}
