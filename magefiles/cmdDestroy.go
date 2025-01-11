package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"
)

func destroy(ctx context.Context) error {
	zap.S().Info("Starting destroy process...")
	config, err := GetConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}

	// Delete Firebase hosting
	if err := cmd.Run(ctx,
		cmd.WithCMD("firebase", "hosting:disable", "--project", config.Env, "--force"),
	); err != nil {
		zap.S().Warnf("failed to disable firebase hosting (might already be disabled): %v", err)
	}

	// Delete Cloud Run service
	if err := cmd.Run(ctx,
		cmd.WithCMD("gcloud", "run", "services", "delete", "gothem-stack",
			"--platform", "managed",
			"--region", "us-central1",
			"--project", config.Env,
			"--quiet"),
	); err != nil {
		zap.S().Warnf("failed to delete cloud run service: %v", err)
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return errors.New("GOOGLE_CLOUD_PROJECT environment variable not set")
	}

	serviceName := os.Getenv("CLOUD_RUN_SERVICE_NAME")
	if serviceName == "" {
		return errors.New("CLOUD_RUN_SERVICE_NAME environment variable not set")
	}

	region := os.Getenv("GOOGLE_CLOUD_REGION")
	if region == "" {
		region = "us-central1"
	}

	// List and delete container images
	err = cmd.Run(ctx,
		cmd.WithCMD("gcloud", "artifacts", "repositories", "delete", serviceName,
			"--project", config.Env,
			"--location", region,
			"--quiet"),
	)
	if err != nil {
		zap.S().Warnf("failed to delete artifacts repository: %v", err)
	}

	// Delete Cloud Build artifacts bucket
	bucketName := fmt.Sprintf("%s_cloudbuild", projectID)
	err = cmd.Run(ctx,
		cmd.WithCMD("gsutil", "rm", "-r", fmt.Sprintf("gs://%s", bucketName)),
	)
	if err != nil {
		zap.S().Warnf("failed to delete cloudbuild bucket: %v", err)
	}

	zap.S().Info("successfully destroyed cloud infrastructure")
	return nil
}
