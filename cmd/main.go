package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/grindlemire/htmx-templ-template/pkg/log"
	"github.com/grindlemire/htmx-templ-template/pkg/server"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	err := log.InitGlobal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initalize logger: %v", err)
		os.Exit(1)
	}

	app := &cli.App{
		Name:  "serve",
		Usage: "serve an htmx api",
		Action: func(c *cli.Context) (err error) {
			ctx, cancel := context.WithCancel(c.Context)
			defer cancel()

			// Create a signal channel
			sigCh := make(chan os.Signal, 1)
			// Register a signal handler for SIGINT
			signal.Notify(sigCh, os.Interrupt)

			go func() {
				<-sigCh
				cancel()
			}()
			return server.Run(ctx)
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		// we don't care about context cancellation as that happens if we kill the process
		// while it is waiting for a request to finish
		if errors.Is(err, context.Canceled) {
			os.Exit(1)
		}
		zap.S().Fatal(err)
	}
}
