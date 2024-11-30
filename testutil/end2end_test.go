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
