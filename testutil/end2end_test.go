//go:build e2e

package testutil

const withTestContainer = false

func TestEnd2End(t *testing.T) {
	var wg sync.WaitGroup
	if withTestContainer {

	} else {
		wg.Add(1)
		go prepareTestEnvironment()
	}

	time.Sleep(5 * time.Second)
	message, err := SentMessage()
	if err != nil {
		t.Fatal(err)
	}
	rabbitmq.GetConnection().Publish("whatsapp.send", message)
	wg.Wait()
	// time.Sleep(5 * time.Second)
	// wg.Done()
}

func prepareTestEnvironment() {
	server.MockStartEnv("../")
	express.Start()
}
