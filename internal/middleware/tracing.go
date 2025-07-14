// Package middleware provides HTTP middleware for the application.
package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"

	"github.com/npmulder/resume-api/internal/tracing"
)

// TracingMiddleware returns a middleware that adds OpenTelemetry tracing to requests.
func TracingMiddleware(tracer *tracing.Tracer) gin.HandlerFunc {
	// Use the otelgin middleware with our configured tracer
	return otelgin.Middleware(
		"resume-api",
		otelgin.WithTracerProvider(tracer.TracerProvider()),
	)
}

// StartSpan starts a new span for the given context and operation name.
// It returns the new context with the span and the span itself.
func StartSpan(ctx *gin.Context, operationName string) (trace.Span, bool) {
	// Extract the tracer from the context
	// This will use the tracer provider that was set by the TracingMiddleware
	tracer := trace.SpanFromContext(ctx.Request.Context()).TracerProvider().Tracer("resume-api")

	// Start a new span
	_, span := tracer.Start(ctx.Request.Context(), operationName)

	// Check if the span is recording (i.e., not a no-op span)
	isRecording := span.IsRecording()

	return span, isRecording
}

// EndSpan ends the given span with the given status.
func EndSpan(span trace.Span, err error) {
	if err != nil {
		// Record error and set status
		span.RecordError(err)
	}

	// End the span
	span.End()
}
