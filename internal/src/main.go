package main

import (
	"github.com/inter-hubly/linker/internal/infraestructure/express"
	"github.com/inter-hubly/pilot/server"
)

func main() {
	// server was start with 'ENVIRONMENT' in tools argument
	server.FillConfigEnvironment()
	express.Start()

	select {}
}
