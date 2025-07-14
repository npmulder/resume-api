package utils

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/npmulder/resume-api/internal/models"
	"github.com/npmulder/resume-api/internal/repository"
)

// ErrorResponse sends a standardized error response to the client
func ErrorResponse(c *gin.Context, status int, message string, opts ...models.APIErrorOption) {
	// Add request path to the error
	pathOpt := models.WithPath(c.Request.URL.Path)
	opts = append(opts, pathOpt)

	// Add request ID if available
	if requestID, exists := c.Get("RequestID"); exists {
		requestIDOpt := models.WithRequestID(requestID.(string))
		opts = append(opts, requestIDOpt)
	}

	// Create the API error
	apiError := models.NewAPIError(status, message, opts...)

	// Send the response
	c.JSON(status, apiError)
	c.Abort()
}

// HandleError handles common error types and returns an appropriate response
func HandleError(c *gin.Context, err error) {
	var repoErr *repository.RepositoryError
	switch {
	case errors.Is(err, repository.ErrNotFound):
		// Handle not found errors
		ErrorResponse(c, http.StatusNotFound, "The requested resource was not found", 
			models.WithCode(models.ErrCodeNotFound))

	case errors.As(err, &repoErr):
		// Handle repository errors
		ErrorResponse(c, http.StatusInternalServerError, "An error occurred while accessing the data",
			models.WithDetails(err.Error()))

	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
		// Handle context errors
		ErrorResponse(c, http.StatusGatewayTimeout, "The request took too long to process",
			models.WithCode(models.ErrCodeServiceUnavailable))

	default:
		// Handle unknown errors
		ErrorResponse(c, http.StatusInternalServerError, "An unexpected error occurred",
			models.WithCode(models.ErrCodeInternalError))
	}
}

// BadRequest returns a bad request error response
func BadRequest(c *gin.Context, message string, details any) {
	opts := []models.APIErrorOption{}

	if details != nil {
		opts = append(opts, models.WithDetails(details))
	}

	ErrorResponse(c, http.StatusBadRequest, message, opts...)
}

// NotFound returns a not found error response
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, models.WithCode(models.ErrCodeNotFound))
}

// InternalError returns an internal server error response
func InternalError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message, models.WithCode(models.ErrCodeInternalError))
}

// ValidationError returns a validation error response
func ValidationError(c *gin.Context, message string, details any) {
	opts := []models.APIErrorOption{models.WithCode(models.ErrCodeValidationFailed)}

	if details != nil {
		opts = append(opts, models.WithDetails(details))
	}

	ErrorResponse(c, http.StatusBadRequest, message, opts...)
}

// Unauthorized returns an unauthorized error response
func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, models.WithCode(models.ErrCodeUnauthorized))
}

// Forbidden returns a forbidden error response
func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, models.WithCode(models.ErrCodeForbidden))
}

// TooManyRequests returns a rate limit exceeded error response
func TooManyRequests(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusTooManyRequests, message, models.WithCode(models.ErrCodeTooManyRequests))
}

// ServiceUnavailable returns a service unavailable error response
func ServiceUnavailable(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusServiceUnavailable, message, models.WithCode(models.ErrCodeServiceUnavailable))
}
