package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"
	"github.com/grindlemire/gothem-stack/pkg/version"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func release(ctx context.Context) error {
	config, err := GetConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}

	// Default to releasing both if no service specified
	if len(config.Args) < 1 {
		config.Args = []string{"all"}
	}

	// Only get version for backend releases
	var ver version.Version
	if config.Args[0] == "backend" || config.Args[0] == "all" {
		ver, err = version.Read()
		if err != nil {
			return errors.Wrap(err, "failed to read version")
		}
	}

	switch config.Args[0] {
	case "backend":
		err = releaseBackend(ctx, ver)
	case "frontend":
		err = releaseFrontend(ctx)
	case "all":
		// Release backend first
		if err = releaseBackend(ctx, ver); err != nil {
			return err
		}
		// Then release frontend
		err = releaseFrontend(ctx)
	default:
		return fmt.Errorf("invalid service: %s. Must be 'backend', 'frontend', or 'all'", config.Args[0])
	}

	if err != nil {
		return err
	}

	if config.Args[0] == "all" {
		zap.S().Info("Successfully released both services")
	} else if config.Args[0] == "backend" {
		zap.S().Infof("Successfully released backend version %s", ver)
	} else {
		zap.S().Infof("Successfully released frontend")
	}
	return nil
}

func releaseBackend(ctx context.Context, ver version.Version) error {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return errors.New("GOOGLE_CLOUD_PROJECT environment variable not set")
	}

	region := os.Getenv("GOOGLE_CLOUD_REGION")
	if region == "" {
		region = "us-central1"
	}

	serviceName := os.Getenv("CLOUD_RUN_SERVICE_NAME")
	if serviceName == "" {
		return errors.New("CLOUD_RUN_SERVICE_NAME environment variable not set")
	}

	// Construct the image tag
	tagname, err := getImageTag(
		projectID,
		serviceName,
		"backend",
		fmt.Sprintf("v%s", ver),
	)
	if err != nil {
		return errors.Wrap(err, "failed to get image tag")
	}

	// Deploy to Cloud Run
	err = cmd.Run(ctx, cmd.WithCMD(
		"gcloud", "run", "deploy", serviceName,
		"--image", tagname,
		"--platform", "managed",
		"--region", region,
		"--project", projectID,
		"--allow-unauthenticated",
	))
	if err != nil {
		return errors.Wrap(err, "failed to deploy to cloud run")
	}

	return nil
}

func releaseFrontend(ctx context.Context) error {
	// if firebase is not installed, error out
	err := cmd.Run(ctx, cmd.WithCMD("firebase", "--version"), cmd.WithSilent())
	if err != nil {
		return errors.Wrap(err, "firebase is not installed. Run `mage install frontend` to install it")
	}

	zap.S().Info("Deploying frontend to Firebase hosting...")

	err = cmd.Run(ctx,
		cmd.WithDir("web"),
		cmd.WithCMD(
			"firebase",
			"deploy",
			"--only", "hosting",
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to deploy to firebase")
	}

	return nil
}
