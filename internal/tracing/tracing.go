// Package tracing provides OpenTelemetry tracing functionality for the application.
package tracing

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/npmulder/resume-api/internal/config"
)

// Tracer is a wrapper around the OpenTelemetry tracer
type Tracer struct {
	tracer trace.Tracer
	tp     *sdktrace.TracerProvider
}

// NewTracer creates a new OpenTelemetry tracer based on the provided configuration
func NewTracer(ctx context.Context, cfg *config.TelemetryConfig, logger *slog.Logger) (*Tracer, error) {
	if !cfg.Enabled {
		logger.Info("tracing is disabled")
		return &Tracer{
			tracer: trace.NewNoopTracerProvider().Tracer("noop"),
		}, nil
	}

	// Create a resource describing the service
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create the appropriate exporter based on the configuration
	var exporter sdktrace.SpanExporter

	switch cfg.ExporterType {
	case "stdout":
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
		}
		logger.Info("using stdout exporter for traces")

	case "otlp":
		// Configure OTLP exporter to send traces to the collector
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.ExporterEndpoint),
			otlptracegrpc.WithInsecure(), // For development; use WithTLSCredentials in production
		}

		client := otlptracegrpc.NewClient(opts...)
		exporter, err = otlptrace.New(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
		}
		logger.Info("using OTLP exporter for traces", "endpoint", cfg.ExporterEndpoint)

	default:
		return nil, fmt.Errorf("unsupported exporter type: %s", cfg.ExporterType)
	}

	// Create a batch span processor for the exporter
	bsp := sdktrace.NewBatchSpanProcessor(exporter)

	// Create a tracer provider with the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SamplingRate)),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Set the global tracer provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Info("tracing initialized", 
		"service", cfg.ServiceName, 
		"exporter", cfg.ExporterType,
		"sampling_rate", cfg.SamplingRate,
	)

	return &Tracer{
		tracer: tp.Tracer(cfg.ServiceName),
		tp:     tp,
	}, nil
}

// Tracer returns the OpenTelemetry tracer
func (t *Tracer) Tracer() trace.Tracer {
	return t.tracer
}

// Shutdown shuts down the tracer provider
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.tp != nil {
		if err := t.tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
	}

	return nil
}
