package main

import (
	"fmt"

	"github.com/pkg/errors"
)

// getRepoName returns the fully qualified repository name for the given service
func getRepoName(projectID, service, image string) string {
	return fmt.Sprintf("us-central1-docker.pkg.dev/%s/%s/gs-%s", projectID, service, image)
}

// getImageTag returns the fully qualified image tag for the given service
func getImageTag(projectID, service, image, version string) (string, error) {
	repoName := getRepoName(projectID, service, image)
	if version == "" {
		return "", errors.New("version cannot be empty")
	}
	return fmt.Sprintf("%s:%s", repoName, version), nil
}
