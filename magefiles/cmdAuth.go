package main

import (
	"context"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// auth authenticates with cloud services
func auth(ctx context.Context) error {
	config, err := GetConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}

	if len(config.Args) == 0 {
		return errors.New("no arguments passed. Please pass 'gcloud' or 'firebase' to authenticate with the respective service")
	}

	switch config.Args[0] {
	case "gcloud":
		return authGcloud(ctx)
	case "firebase":
		return authFirebase(ctx)
	default:
		return errors.Errorf("unknown auth target: %s", config.Args[0])
	}
}

func authGcloud(ctx context.Context) error {
	zap.S().Info("Authenticating with Google Cloud...")

	// Check if gcloud is installed
	err := cmd.Run(ctx, cmd.WithCMD("gcloud", "version"), cmd.WithSilent())
	if err != nil {
		return errors.Wrap(err, "gcloud is not installed. Run `mage install deploy` to install it")
	}

	// Run gcloud auth login
	err = cmd.Run(ctx, cmd.WithCMD(
		"gcloud", "auth", "login",
	))
	if err != nil {
		return errors.Wrap(err, "failed to authenticate with gcloud")
	}

	zap.S().Info("Successfully authenticated with Google Cloud")
	return nil
}

func authFirebase(ctx context.Context) error {
	zap.S().Info("Authenticating with Firebase...")

	// Check if firebase is installed
	err := cmd.Run(ctx, cmd.WithCMD("firebase", "--version"), cmd.WithSilent())
	if err != nil {
		return errors.Wrap(err, "firebase is not installed. Run `mage install deploy` to install it")
	}

	// Run firebase login
	err = cmd.Run(ctx, cmd.WithCMD(
		"firebase", "login",
	))
	if err != nil {
		return errors.Wrap(err, "failed to authenticate with firebase")
	}

	zap.S().Info("Successfully authenticated with Firebase")
	return nil
}
