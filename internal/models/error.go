package models

import (
	"fmt"
	"net/http"
	"time"
)

// APIError represents a standardized API error response
type APIError struct {
	Status     int       `json:"status"`               // HTTP status code
	Code       string    `json:"code"`                 // Application-specific error code
	Message    string    `json:"message"`              // Human-readable error message
	Details    any       `json:"details,omitempty"`    // Additional error details (optional)
	RequestID  string    `json:"request_id,omitempty"` // Request ID for tracing (optional)
	Timestamp  time.Time `json:"timestamp"`            // Time when the error occurred
	Path       string    `json:"path,omitempty"`       // Request path that caused the error
	Suggestion string    `json:"suggestion,omitempty"` // Suggested action to resolve the error (optional)
}

// Error codes
const (
	// General errors
	ErrCodeInternalError     = "INTERNAL_ERROR"
	ErrCodeBadRequest        = "BAD_REQUEST"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeValidationFailed  = "VALIDATION_FAILED"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeTooManyRequests   = "TOO_MANY_REQUESTS"
	ErrCodeMethodNotAllowed  = "METHOD_NOT_ALLOWED"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	
	// Resource-specific errors
	ErrCodeProfileNotFound   = "PROFILE_NOT_FOUND"
	ErrCodeExperienceNotFound = "EXPERIENCE_NOT_FOUND"
	ErrCodeSkillNotFound     = "SKILL_NOT_FOUND"
	ErrCodeEducationNotFound = "EDUCATION_NOT_FOUND"
	ErrCodeProjectNotFound   = "PROJECT_NOT_FOUND"
	ErrCodeAchievementNotFound = "ACHIEVEMENT_NOT_FOUND"
)

// HTTP status code to error code mapping
var statusToErrorCode = map[int]string{
	http.StatusBadRequest:          ErrCodeBadRequest,
	http.StatusUnauthorized:        ErrCodeUnauthorized,
	http.StatusForbidden:           ErrCodeForbidden,
	http.StatusNotFound:            ErrCodeNotFound,
	http.StatusMethodNotAllowed:    ErrCodeMethodNotAllowed,
	http.StatusInternalServerError: ErrCodeInternalError,
	http.StatusServiceUnavailable:  ErrCodeServiceUnavailable,
	http.StatusTooManyRequests:     ErrCodeTooManyRequests,
}

// GetErrorCodeForStatus returns the appropriate error code for a given HTTP status
func GetErrorCodeForStatus(status int) string {
	if code, exists := statusToErrorCode[status]; exists {
		return code
	}
	return ErrCodeInternalError
}

// NewAPIError creates a new APIError with the given parameters
func NewAPIError(status int, message string, opts ...APIErrorOption) *APIError {
	err := &APIError{
		Status:    status,
		Code:      GetErrorCodeForStatus(status),
		Message:   message,
		Timestamp: time.Now().UTC(),
	}

	// Apply all options
	for _, opt := range opts {
		opt(err)
	}

	return err
}

// APIErrorOption is a function that configures an APIError
type APIErrorOption func(*APIError)

// WithCode sets a custom error code
func WithCode(code string) APIErrorOption {
	return func(e *APIError) {
		e.Code = code
	}
}

// WithDetails adds additional error details
func WithDetails(details any) APIErrorOption {
	return func(e *APIError) {
		e.Details = details
	}
}

// WithRequestID sets the request ID for tracing
func WithRequestID(requestID string) APIErrorOption {
	return func(e *APIError) {
		e.RequestID = requestID
	}
}

// WithPath sets the request path
func WithPath(path string) APIErrorOption {
	return func(e *APIError) {
		e.Path = path
	}
}

// WithSuggestion adds a suggestion for resolving the error
func WithSuggestion(suggestion string) APIErrorOption {
	return func(e *APIError) {
		e.Suggestion = suggestion
	}
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s (status: %d)", e.Code, e.Message, e.Status)
}