package main

import (
	"context"

	"github.com/grindlemire/htmx-templ-template/magefiles/cmd"

	"go.uber.org/zap"
)

func tidy(ctx context.Context) (err error) {
	err = cmd.Run(ctx,
		cmd.WithCMD(
			"go",
			"mod", "tidy",
		),
	)
	if err != nil {
		zap.S().Errorf("Error running go mod tidy: %v", err)
	}
	return err

}
