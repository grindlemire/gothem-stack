package main

import (
	"context"
	"runtime"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"
	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

// backend installs backend-specific dependencies
func backend(ctx context.Context) error {
	zap.S().Info("Installing backend dependencies...")

	// Install gcloud CLI
	if runtime.GOOS == "darwin" {
		err := cmd.Run(ctx, cmd.WithCMD(
			"brew", "install", "google-cloud-sdk",
		))
		if err != nil {
			return err
		}
	} else if runtime.GOOS == "linux" {
		return errors.New("Please install the gcloud CLI manually")
	}

	// Install Firebase CLI globally
	err := cmd.Run(ctx, cmd.WithCMD(
		"npm", "install", "-g", "firebase-tools",
	))
	if err != nil {
		return err
	}

	zap.S().Info("Deployment dependencies installed successfully")
	return nil
}

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

	zap.S().Info("installing frontend dependencies")
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

	// install backend dependencies if the backend arg is passed
	if config.Args[0] == "deploy" {
		err = backend(ctx)
		if err != nil {
			return err
		}
	}

	mg.SerialCtxDeps(ctx, tidy)
	return err
}
