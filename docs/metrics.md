# Metrics and Observability

This document describes the metrics and observability features implemented in the Resume API.

## Overview

The Resume API uses OpenTelemetry for metrics collection with a Prometheus exporter for compatibility. Metrics are collected for:

- HTTP requests (count, duration, in-flight)
- Database operations (count, duration)
- System resources (memory usage, goroutines count)

## OpenTelemetry Integration

OpenTelemetry is an observability framework for cloud-native software that provides a collection of tools, APIs, and SDKs for distributed tracing, metrics, and logging. The Resume API uses OpenTelemetry for metrics collection with the following components:

- OpenTelemetry SDK for Go
- Prometheus exporter for compatibility with existing monitoring systems
- Automatic instrumentation of HTTP requests and database operations

## Metrics Endpoint

Metrics are exposed at the `/metrics` endpoint in Prometheus format. This endpoint can be scraped by a Prometheus server to collect and store the metrics.

## Available Metrics

### HTTP Metrics

- `http_requests_total` - Total number of HTTP requests by method, path, and status
- `http_request_duration_seconds` - Duration of HTTP requests in seconds
- `http_requests_in_flight` - Current number of HTTP requests in flight

### Database Metrics

- `database_operations_total` - Total number of database operations by operation type
- `database_operation_duration_seconds` - Duration of database operations in seconds

### System Metrics

- `memory_usage_bytes` - Current memory usage in bytes (alloc, sys, heap_alloc, heap_sys)
- `goroutines_count` - Current number of goroutines

## Using Database Operation Tracking

To track database operations in your repository implementations, use the `TrackDatabaseOperation` function:

```go
import "github.com/npmulder/resume-api/internal/middleware"

func (r *profileRepository) GetProfile(ctx context.Context) (*models.Profile, error) {
    var profile models.Profile

    err := middleware.TrackDatabaseOperation("get_profile", func() error {
        // Database operation here
        return r.db.QueryRow(ctx, "SELECT * FROM profiles LIMIT 1").Scan(&profile.ID, &profile.Name, ...)
    })

    if err != nil {
        return nil, err
    }

    return &profile, nil
}
```

## Prometheus Configuration

To scrape these metrics with Prometheus, add the following to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: 'resume-api'
    scrape_interval: 15s
    static_configs:
      - targets: ['resume-api:8080']
```

## Grafana Dashboard

A sample Grafana dashboard can be created to visualize these metrics. Here are some useful panels to include:

- Request Rate: `rate(http_requests_total[1m])`
- Request Duration: `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))`
- Error Rate: `sum(rate(http_requests_total{status=~"5.."}[1m])) / sum(rate(http_requests_total[1m]))`
- Memory Usage: `memory_usage_bytes{type="heap_alloc"}`
- Goroutines Count: `goroutines_count`

## Future Improvements

- Add CPU usage metrics
- Add connection pool metrics
- Add custom business metrics
- Implement distributed tracing with OpenTelemetry's tracing API
