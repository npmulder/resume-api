// Package versioning provides utilities for API version management.
package versioning

import (
	"fmt"
	"strings"
)

// Version represents an API version
type Version string

// Supported API versions
const (
	// V1 is the initial API version
	V1 Version = "v1"
	
	// Add new versions here as they are developed
	// V2 Version = "v2"
	
	// LatestVersion should always point to the most recent stable version
	LatestVersion = V1
)

// All returns all supported API versions
func All() []Version {
	return []Version{V1}
}

// IsSupported checks if the given version is supported
func IsSupported(v string) bool {
	v = strings.TrimPrefix(strings.ToLower(v), "v")
	for _, version := range All() {
		if strings.TrimPrefix(string(version), "v") == v {
			return true
		}
	}
	return false
}

// Normalize ensures the version is in the standard format (e.g., "v1")
func Normalize(v string) (Version, error) {
	// Remove any "v" prefix and convert to lowercase
	v = strings.TrimPrefix(strings.ToLower(v), "v")
	
	// Add "v" prefix back
	versionStr := "v" + v
	
	// Check if it's a supported version
	if !IsSupported(versionStr) {
		return "", fmt.Errorf("unsupported API version: %s", versionStr)
	}
	
	return Version(versionStr), nil
}

// GetPathPrefix returns the API path prefix for a given version (e.g., "/api/v1")
func GetPathPrefix(v Version) string {
	return fmt.Sprintf("/api/%s", v)
}

// GetLatestPathPrefix returns the API path prefix for the latest version
func GetLatestPathPrefix() string {
	return GetPathPrefix(LatestVersion)
}