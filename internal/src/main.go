package main

import (
	"context"

	"github.com/inter-hubly/linker/internal/infraestructure/express"
	"github.com/inter-hubly/pilot/server"
)

func main() {
	// server was start with 'ENVIRONMENT' in tools argument
	ctx := context.Background()
	server.FillConfigEnvironment(ctx)
	express.Start(ctx)

	select {}
}
