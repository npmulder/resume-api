package middleware

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var (
	// OpenTelemetry meter provider and meter
	meterProvider *sdkmetric.MeterProvider
	meter         metric.Meter

	// HTTP metrics
	httpRequestsTotal       metric.Int64Counter
	httpRequestDuration     metric.Float64Histogram
	httpRequestsInFlight    metric.Int64UpDownCounter

	// Database metrics
	dbOperationsTotal       metric.Int64Counter
	dbOperationDuration     metric.Float64Histogram

	// System metrics
	memoryUsage             metric.Float64ObservableGauge
	goroutinesCount         metric.Int64ObservableGauge

	// Initialization flag
	initialized             bool
	initMutex               sync.Mutex
)

// initMetrics initializes the OpenTelemetry metrics
func initMetrics() error {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return nil
	}

	// Create a Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create a meter provider with the Prometheus exporter
	meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)

	// Set the global meter provider
	otel.SetMeterProvider(meterProvider)

	// Create a meter
	meter = meterProvider.Meter("github.com/npmulder/resume-api")

	// Create HTTP metrics
	httpRequestsTotal, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http_requests_total counter: %w", err)
	}

	httpRequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http_request_duration_seconds histogram: %w", err)
	}

	httpRequestsInFlight, err = meter.Int64UpDownCounter(
		"http_requests_in_flight",
		metric.WithDescription("Current number of HTTP requests in flight"),
	)
	if err != nil {
		return fmt.Errorf("failed to create http_requests_in_flight counter: %w", err)
	}

	// Create database metrics
	dbOperationsTotal, err = meter.Int64Counter(
		"database_operations_total",
		metric.WithDescription("Total number of database operations"),
	)
	if err != nil {
		return fmt.Errorf("failed to create database_operations_total counter: %w", err)
	}

	dbOperationDuration, err = meter.Float64Histogram(
		"database_operation_duration_seconds",
		metric.WithDescription("Duration of database operations in seconds"),
	)
	if err != nil {
		return fmt.Errorf("failed to create database_operation_duration_seconds histogram: %w", err)
	}

	// Create system metrics
	memoryUsage, err = meter.Float64ObservableGauge(
		"memory_usage_bytes",
		metric.WithDescription("Current memory usage in bytes"),
	)
	if err != nil {
		return fmt.Errorf("failed to create memory_usage_bytes gauge: %w", err)
	}

	goroutinesCount, err = meter.Int64ObservableGauge(
		"goroutines_count",
		metric.WithDescription("Current number of goroutines"),
	)
	if err != nil {
		return fmt.Errorf("failed to create goroutines_count gauge: %w", err)
	}

	// Register callbacks for observable metrics
	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			o.ObserveFloat64(memoryUsage, float64(memStats.Alloc), metric.WithAttributes(attribute.String("type", "alloc")))
			o.ObserveFloat64(memoryUsage, float64(memStats.Sys), metric.WithAttributes(attribute.String("type", "sys")))
			o.ObserveFloat64(memoryUsage, float64(memStats.HeapAlloc), metric.WithAttributes(attribute.String("type", "heap_alloc")))
			o.ObserveFloat64(memoryUsage, float64(memStats.HeapSys), metric.WithAttributes(attribute.String("type", "heap_sys")))

			o.ObserveInt64(goroutinesCount, int64(runtime.NumGoroutine()))

			return nil
		},
		memoryUsage,
		goroutinesCount,
	)
	if err != nil {
		return fmt.Errorf("failed to register callback: %w", err)
	}

	initialized = true
	return nil
}

// MetricsMiddleware returns a middleware that collects HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	// Initialize metrics
	if err := initMetrics(); err != nil {
		panic(fmt.Sprintf("failed to initialize metrics: %v", err))
	}

	return func(c *gin.Context) {
		// Skip metrics endpoint to avoid circular measurements
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Track in-flight requests
		ctx := c.Request.Context()
		httpRequestsInFlight.Add(ctx, 1)
		defer httpRequestsInFlight.Add(ctx, -1)

		// Record start time
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// Create attributes for the metrics
		attrs := []attribute.KeyValue{
			attribute.String("method", c.Request.Method),
			attribute.String("path", path),
			attribute.String("status", status),
		}

		httpRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
		httpRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
	}
}

// TrackDatabaseOperation is a utility function to track database operations
func TrackDatabaseOperation(operation string, f func() error) error {
	// Initialize metrics if not already initialized
	if err := initMetrics(); err != nil {
		// Log the error but don't fail the operation
		fmt.Printf("failed to initialize metrics: %v\n", err)
		return f()
	}

	ctx := context.Background()
	start := time.Now()
	err := f()
	duration := time.Since(start).Seconds()

	// Create attributes for the metrics
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
	}

	dbOperationsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	dbOperationDuration.Record(ctx, duration, metric.WithAttributes(attrs...))

	return err
}
