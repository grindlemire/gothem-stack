package main

import (
	"context"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func deploy(ctx context.Context) error {
	zap.S().Info("Deploying to Firebase hosting...")

	// First, ensure we have the latest static build
	if err := static(ctx); err != nil {
		return errors.Wrap(err, "failed to build static files")
	}

	// Deploy to Firebase using local node_modules installation
	err := cmd.Run(ctx,
		cmd.WithDir("web"),
		cmd.WithCMD(
			"./node_modules/.bin/firebase",
			"deploy",
			"--only", "hosting",
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to deploy to firebase")
	}

	zap.S().Info("Successfully deployed to Firebase hosting")
	return nil
}
