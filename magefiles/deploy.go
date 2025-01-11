package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func deploy(ctx context.Context) error {
	zap.S().Info("Starting deployment process...")
	config, err := GetConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}

	if len(config.Args) == 0 {
		return errors.New("no arguments passed. Please pass 'backend' or 'frontend' to deploy the respective services")
	}

	// First, ensure we have the latest static build
	if err := static(ctx); err != nil {
		return errors.Wrap(err, "failed to build static files")
	}

	if config.Args[0] == "backend" || config.Args[0] == "all" {
		err = deployBackend(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to deploy backend to Cloud Run")
		}
	}

	if config.Args[0] == "frontend" || config.Args[0] == "all" {
		err = deployFrontend(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to deploy to firebase")
		}
	}

	zap.S().Info("Successfully completed deployment")
	return nil
}

func deployBackend(ctx context.Context) error {
	zap.S().Info("Deploying backend to Cloud Run...")
	// check if gcloud is installed
	err := cmd.Run(ctx, cmd.WithCMD("gcloud", "version"))
	if err != nil {
		return errors.Wrap(err, "gcloud is not installed. Run `mage install backend` to install it")
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return errors.New("GOOGLE_CLOUD_PROJECT environment variable not set")
	}

	region := os.Getenv("GOOGLE_CLOUD_REGION")
	if region == "" {
		region = "us-central1" // Default region
	}

	serviceName := "gothem-backend"
	imageTag := fmt.Sprintf("gcr.io/%s/%s:latest", projectID, serviceName)

	// Build and push using Cloud Build
	err = cmd.Run(ctx, cmd.WithCMD(
		"gcloud", "builds", "submit",
		"--tag", imageTag,
	))
	if err != nil {
		return errors.Wrap(err, "failed to build and push image")
	}

	// Deploy to Cloud Run
	err = cmd.Run(ctx, cmd.WithCMD(
		"gcloud", "run", "deploy", serviceName,
		"--image", imageTag,
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

func deployFrontend(ctx context.Context) error {
	// if firebase is not installed, error out
	err := cmd.Run(ctx, cmd.WithCMD("firebase", "version"))
	if err != nil {
		return errors.Wrap(err, "firebase is not installed. Run `mage install frontend` to install it")
	}

	zap.S().Info("Deploying to Firebase hosting...")

	// Deploy to Firebase using local node_modules installation
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

	zap.S().Info("Successfully deployed to Firebase hosting")
	return nil
}
