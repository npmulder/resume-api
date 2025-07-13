// Package versioning provides utilities for API version management.
package versioning

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// VersionNegotiationOptions configures how version negotiation works
type VersionNegotiationOptions struct {
	// EnableURIPath enables version detection from URI path (e.g., /api/v1/profile)
	EnableURIPath bool
	
	// EnableAcceptHeader enables version detection from Accept header
	// (e.g., Accept: application/json;version=1)
	EnableAcceptHeader bool
	
	// EnableQueryParam enables version detection from query parameter
	// (e.g., ?version=1)
	EnableQueryParam bool
	
	// QueryParamName is the name of the query parameter for version detection
	QueryParamName string
	
	// DefaultToLatest determines if requests without a version should use the latest version
	DefaultToLatest bool
}

// DefaultVersionNegotiationOptions returns the default options for version negotiation
func DefaultVersionNegotiationOptions() VersionNegotiationOptions {
	return VersionNegotiationOptions{
		EnableURIPath:     true,
		EnableAcceptHeader: true,
		EnableQueryParam:  true,
		QueryParamName:    "version",
		DefaultToLatest:   true,
	}
}

// VersionKey is the key used to store the API version in the Gin context
const VersionKey = "api_version"

// VersionNegotiationMiddleware creates middleware that determines the requested API version
// It tries multiple methods based on the provided options and sets the version in the context
func VersionNegotiationMiddleware(options VersionNegotiationOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var version Version
		var found bool
		
		// Try to extract version from URI path
		if !found && options.EnableURIPath {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api/") {
				parts := strings.Split(path, "/")
				if len(parts) >= 3 {
					potentialVersion := parts[2]
					if v, err := Normalize(potentialVersion); err == nil {
						version = v
						found = true
					}
				}
			}
		}
		
		// Try to extract version from Accept header
		if !found && options.EnableAcceptHeader {
			accept := c.GetHeader("Accept")
			if strings.Contains(accept, "version=") {
				parts := strings.Split(accept, "version=")
				if len(parts) >= 2 {
					versionPart := parts[1]
					endIndex := strings.IndexAny(versionPart, ";,")
					if endIndex != -1 {
						versionPart = versionPart[:endIndex]
					}
					if v, err := Normalize(versionPart); err == nil {
						version = v
						found = true
					}
				}
			}
		}
		
		// Try to extract version from query parameter
		if !found && options.EnableQueryParam {
			queryVersion := c.Query(options.QueryParamName)
			if queryVersion != "" {
				if v, err := Normalize(queryVersion); err == nil {
					version = v
					found = true
				}
			}
		}
		
		// If no version found and DefaultToLatest is true, use the latest version
		if !found && options.DefaultToLatest {
			version = LatestVersion
			found = true
		}
		
		// If a version was found, set it in the context
		if found {
			c.Set(VersionKey, version)
			c.Next()
		} else {
			// No valid version found
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Unsupported API version",
			})
			c.Abort()
		}
	}
}

// GetRequestedVersion retrieves the API version from the Gin context
func GetRequestedVersion(c *gin.Context) Version {
	if v, exists := c.Get(VersionKey); exists {
		if version, ok := v.(Version); ok {
			return version
		}
	}
	return LatestVersion
}