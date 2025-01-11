package version

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// DefaultVersionFile is the default path for the version file
const DefaultVersionFile = "version.txt"

// Version represents the version of a service
type Version struct {
	Major int
	Minor int
	Patch int
}

// String returns the version as a string
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// ReadFromFile reads the version from the specified file path
func ReadFromFile(filepath string) (Version, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			// Start with 0.0.1 if file doesn't exist
			return Version{Major: 0, Minor: 0, Patch: 1}, nil
		}
		return Version{}, errors.Wrap(err, "failed to read version file")
	}
	return parseVersion(strings.TrimSpace(string(data)))
}

// WriteToFile writes the version to the specified file path
func WriteToFile(v Version, filepath string) error {
	return os.WriteFile(filepath, []byte(v.String()), 0644)
}

// Read reads the version from the default version file
func Read() (Version, error) {
	return ReadFromFile(DefaultVersionFile)
}

// Write writes the version to the default version file
func Write(v Version) error {
	return WriteToFile(v, DefaultVersionFile)
}

// Increment increments the version
func (v *Version) Increment(level string) error {
	switch level {
	case "major":
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case "minor":
		v.Minor++
		v.Patch = 0
	case "patch":
		v.Patch++
	default:
		return fmt.Errorf("invalid version increment level: %s", level)
	}
	return nil
}

// parseVersion parses the version from a string
func parseVersion(ver string) (Version, error) {
	parts := strings.Split(strings.TrimPrefix(ver, "v"), ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("invalid version format: %s", ver)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, errors.Wrap(err, "invalid major version")
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, errors.Wrap(err, "invalid minor version")
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, errors.Wrap(err, "invalid patch version")
	}

	return Version{Major: major, Minor: minor, Patch: patch}, nil
}
