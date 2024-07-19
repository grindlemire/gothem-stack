package main

import (
	"context"

	"github.com/grindlemire/htmx-templ-template/magefiles/cmd"
	"github.com/magefile/mage/mg"

	"go.uber.org/zap"
)

func build(ctx context.Context) error {
	config, err := GetConfig(ctx)
	if err != nil {
		return err
	}

	zap.S().Infof("mage running with configs: %+v", config)

	mg.SerialCtxDeps(ctx, tidy)

	err = cmd.Run(ctx,
		cmd.WithCMD(
			"rm", "-f", "dist/server",
		),
	)
	if err != nil {
		return err
	}

	err = cmd.Run(ctx,
		cmd.WithCMD(
			"go",
			"build",
			"-o", "dist/server",
			"cmd/main.go",
		),
	)
	if err != nil {
		zap.S().Errorf("Error running build command: %v", err)
	}
	return err
}
