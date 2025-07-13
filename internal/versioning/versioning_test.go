package versioning

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestVersionNegotiationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		path           string
		acceptHeader   string
		queryParam     string
		options        VersionNegotiationOptions
		expectedStatus int
		expectedVersion Version
	}{
		{
			name:           "URI Path Version",
			path:           "/api/v1/profile",
			options:        DefaultVersionNegotiationOptions(),
			expectedStatus: http.StatusOK,
			expectedVersion: V1,
		},
		{
			name:           "Accept Header Version",
			path:           "/api/profile",
			acceptHeader:   "application/json;version=1",
			options:        DefaultVersionNegotiationOptions(),
			expectedStatus: http.StatusOK,
			expectedVersion: V1,
		},
		{
			name:           "Query Param Version",
			path:           "/api/profile",
			queryParam:     "version=1",
			options:        DefaultVersionNegotiationOptions(),
			expectedStatus: http.StatusOK,
			expectedVersion: V1,
		},
		{
			name:           "Default to Latest Version",
			path:           "/api/profile",
			options:        DefaultVersionNegotiationOptions(),
			expectedStatus: http.StatusOK,
			expectedVersion: LatestVersion,
		},
		{
			name:           "Unsupported Version",
			path:           "/api/v999/profile",
			options:        VersionNegotiationOptions{
				EnableURIPath:     true,
				EnableAcceptHeader: false,
				EnableQueryParam:  false,
				DefaultToLatest:   false,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(VersionNegotiationMiddleware(tt.options))
			
			// Add a test handler that returns the version from context
			router.GET("/*path", func(c *gin.Context) {
				version := GetRequestedVersion(c)
				c.JSON(http.StatusOK, gin.H{"version": version})
			})
			
			// Create test request
			req, _ := http.NewRequest("GET", tt.path, nil)
			if tt.acceptHeader != "" {
				req.Header.Set("Accept", tt.acceptHeader)
			}
			if tt.queryParam != "" {
				req.URL.RawQuery = tt.queryParam
			}
			
			// Perform the request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			// If we expect success, check the version
			if tt.expectedStatus == http.StatusOK {
				assert.Contains(t, w.Body.String(), string(tt.expectedVersion))
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      Version
		expectError   bool
	}{
		{
			name:        "Valid version with v prefix",
			input:       "v1",
			expected:    V1,
			expectError: false,
		},
		{
			name:        "Valid version without v prefix",
			input:       "1",
			expected:    V1,
			expectError: false,
		},
		{
			name:        "Valid version with uppercase",
			input:       "V1",
			expected:    V1,
			expectError: false,
		},
		{
			name:        "Invalid version",
			input:       "v999",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Normalize(tt.input)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}