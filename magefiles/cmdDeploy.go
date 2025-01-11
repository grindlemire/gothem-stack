package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/grindlemire/gothem-stack/magefiles/cmd"
	"github.com/grindlemire/gothem-stack/pkg/version"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func deploy(ctx context.Context) error {
	zap.S().Info("Starting deployment process...")
	config, err := GetConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}

	// If no arguments are passed, deploy both backend and frontend
	if len(config.Args) == 0 {
		config.Args = []string{"all"}
	}

	var ver version.Version
	// Only handle versioning for backend
	if config.Args[0] == "backend" || config.Args[0] == "all" {
		// Default to patch if no version increment specified
		versionLevel := "patch"
		if len(config.Args) > 1 {
			versionLevel = config.Args[1]
		}

		// Read current version
		ver, err = version.Read()
		if err != nil {
			return errors.Wrap(err, "failed to read version")
		}

		// Increment version
		if err := ver.Increment(versionLevel); err != nil {
			return errors.Wrap(err, "failed to increment version")
		}
	}

	// First, ensure we have the latest build
	if err := build(ctx); err != nil {
		return errors.Wrap(err, "failed to build static files")
	}

	if config.Args[0] == "backend" || config.Args[0] == "all" {
		err = deployBackend(ctx, ver)
		if err != nil {
			return errors.Wrap(err, "failed to deploy backend to Cloud Run")
		}
		// Save new version after successful backend deployment
		if err := version.Write(ver); err != nil {
			return errors.Wrap(err, "failed to save new version")
		}
		zap.S().Infof("Successfully deployed backend version %s", ver)
	}

	if config.Args[0] == "frontend" || config.Args[0] == "all" {
		err = deployFrontend(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to deploy to firebase")
		}
		zap.S().Info("Successfully deployed frontend")
	}

	zap.S().Info("Successfully completed deployment")
	return nil
}

func deployBackend(ctx context.Context, ver version.Version) error {
	zap.S().Info("Deploying backend to Cloud Run...")
	// check if gcloud is installed
	err := cmd.Run(ctx, cmd.WithCMD("gcloud", "version"), cmd.WithSilent())
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

	serviceName := os.Getenv("CLOUD_RUN_SERVICE_NAME")
	if serviceName == "" {
		return errors.New("CLOUD_RUN_SERVICE_NAME environment variable not set")
	}

	// Ensure the Artifact Registry is initialized
	err = ensureArtifactRegistry(ctx, projectID, region, serviceName)
	if err != nil {
		return errors.Wrap(err, "failed to initialize artifact registry")
	}

	err = ensureCloudBuild(ctx, projectID, region, serviceName)
	if err != nil {
		return errors.Wrap(err, "failed to initialize cloud build")
	}

	err = ensureCloudRun(ctx, projectID)
	if err != nil {
		return errors.Wrap(err, "failed to initialize cloud run")
	}

	// create new tag name in accordance to the format gcr requires
	tagname, err := getImageTag(
		projectID,
		serviceName,
		"backend",
		fmt.Sprintf("v%s", ver),
	)
	if err != nil {
		return errors.Wrap(err, "failed to get image tag")
	}

	err = cmd.Run(ctx,
		cmd.WithCMD(
			"gcloud", "builds", "submit",
			"--project", projectID,
			"--region", region,
			"--config", "./cloudbuild.json",
			"--substitutions", fmt.Sprintf("_IMAGE_NAME=%s", tagname),
		),
	)
	if err != nil {
		return err
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

	zap.S().Info("Successfully deployed to Cloud Run")
	return nil
}

func ensureCloudRun(ctx context.Context, projectID string) error {
	// Check if Cloud Run API is enabled
	err := cmd.Run(ctx, cmd.WithCMD(
		"gcloud", "services", "enable", "run.googleapis.com",
		"--project", projectID,
	))
	if err != nil {
		return errors.Wrap(err, "failed to enable Cloud Run API")
	}
	return nil
}

func ensureArtifactRegistry(ctx context.Context, projectID, region, repoName string) error {
	// Check if Artifact Registry API is enabled
	err := cmd.Run(ctx, cmd.WithCMD(
		"gcloud", "services", "enable", "artifactregistry.googleapis.com",
		"--project", projectID,
	))
	if err != nil {
		return errors.Wrap(err, "failed to enable Artifact Registry API")
	}

	// Check if repository exists
	repoExists := cmd.Run(ctx,
		cmd.WithCMD(
			"gcloud", "artifacts", "repositories", "describe", repoName,
			"--project", projectID,
			"--location", region,
		),
		cmd.WithSilent(),
	) == nil

	// Create repository if it doesn't exist
	if !repoExists {
		zap.S().Infof("Creating Artifact Registry repository %s in %s", repoName, region)
		err = cmd.Run(ctx, cmd.WithCMD(
			"gcloud", "artifacts", "repositories", "create", repoName,
			"--repository-format", "docker",
			"--location", region,
			"--project", projectID,
		))
		if err != nil {
			return errors.Wrap(err, "failed to create Artifact Registry repository")
		}
	}

	// Set cleanup policy
	zap.S().Infof("Setting cleanup policy for repository %s", repoName)
	err = cmd.Run(ctx, cmd.WithCMD(
		"gcloud", "artifacts", "repositories", "set-cleanup-policies", repoName,
		"--location", region,
		"--project", projectID,
		"--policy=artifact-cleanup-policy.json",
	))
	if err != nil {
		return errors.Wrap(err, "failed to set cleanup policy")
	}

	return nil
}

func ensureCloudBuild(ctx context.Context, projectID, region, serviceName string) error {
	// Enable Cloud Build API
	err := cmd.Run(ctx, cmd.WithCMD(
		"gcloud", "services", "enable", "cloudbuild.googleapis.com",
		"--project", projectID,
	))
	if err != nil {
		return errors.Wrap(err, "failed to enable Cloud Build API")
	}
	return nil
}

func deployFrontend(ctx context.Context) error {
	// Removed version parameter
	// if firebase is not installed, error out
	err := cmd.Run(ctx, cmd.WithCMD("firebase", "--version"), cmd.WithSilent())
	if err != nil {
		return errors.Wrap(err, "firebase is not installed. Run `mage install frontend` to install it")
	}

	zap.S().Info("Deploying to Firebase hosting...")

	err = cmd.Run(ctx,
		cmd.WithDir("web"),
		cmd.WithCMD(
			"firebase",
			"deploy",
			"--only", "hosting",
			"--message", fmt.Sprintf("Build at %s", time.Now().Format(time.RFC3339)),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to deploy to firebase")
	}

	zap.S().Info("Successfully deployed to Firebase hosting")
	return nil
}
