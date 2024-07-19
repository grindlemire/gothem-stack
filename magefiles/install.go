package main

import (
	"context"

	"github.com/grindlemire/htmx-templ-template/magefiles/cmd"
	"github.com/magefile/mage/mg"

	"go.uber.org/zap"
)

func install(ctx context.Context) error {
	config, err := GetConfig(ctx)
	if err != nil {
		return err
	}

	zap.S().Infof("mage running with configs: %+v", config)

	// pin to a specific commit for now. See https://github.com/air-verse/air/issues/534
	zap.S().Info("installing air at pinned commit for #534")
	err = cmd.Run(ctx,
		cmd.WithCMD(
			"go",
			"install",
			"github.com/air-verse/air@360714a021b1b77e50a5656fefc4f8bb9312d328",
		),
	)
	if err != nil {
		return err
	}

	// pin to a specific commit for now. See https://github.com/a-h/templ/pull/841
	zap.S().Info("installing templ at pinned commit for #841")
	err = cmd.Run(ctx,
		cmd.WithCMD(
			"go",
			"install",
			"github.com/a-h/templ/cmd/templ@v0.2.707",
		),
	)
	if err != nil {
		return err
	}

	zap.S().Info("running go mod tidy")
	mg.SerialCtxDeps(ctx, templ, tidy)

	zap.S().Info("installing frontend deps")
	err = cmd.Run(ctx,
		cmd.WithDir("web"),
		cmd.WithCMD(
			"npm",
			"install",
		),
	)
	if err != nil {
		return err
	}

	mg.SerialCtxDeps(ctx, tidy)
	return err
}
