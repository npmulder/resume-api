# Distributed Tracing with OpenTelemetry

This document describes the distributed tracing implementation in the Resume API using OpenTelemetry.

## Overview

The Resume API uses OpenTelemetry for distributed tracing, which provides insights into the application's behavior and performance. Tracing helps identify bottlenecks, understand request flows, and debug issues in production.

## Configuration

Tracing is configured through the application's configuration system. The following settings are available:

```yaml
telemetry:
  enabled: true                # Enable or disable tracing
  service_name: "resume-api"   # Name of the service in traces
  exporter_type: "stdout"      # Exporter type: stdout, otlp
  exporter_endpoint: ""        # Endpoint for the exporter (not needed for stdout)
  sampling_rate: 1.0           # Sampling rate (0.0 to 1.0)
```

These settings can be configured through environment variables:

```
RESUME_API_TELEMETRY_ENABLED=true
RESUME_API_TELEMETRY_SERVICE_NAME=resume-api
RESUME_API_TELEMETRY_EXPORTER_TYPE=stdout
RESUME_API_TELEMETRY_EXPORTER_ENDPOINT=
RESUME_API_TELEMETRY_SAMPLING_RATE=1.0
```

## Exporters

The following exporters are supported:

- **stdout**: Prints traces to standard output (useful for development)
- **otlp**: Sends traces to an OpenTelemetry collector using OTLP protocol

## Instrumentation

The Resume API has instrumentation at multiple levels:

### HTTP Layer

All HTTP requests are automatically traced using the Gin middleware. The middleware captures:

- Request method and path
- Response status code
- Request duration
- Error information (if any)

### Database Layer

Database operations are traced using a custom wrapper around the pgx connection pool. The wrapper captures:

- SQL queries (without parameters for security)
- Query duration
- Error information (if any)
- Number of rows affected (for write operations)

## Usage in Code

### Starting a New Span

To create a new span in your code:

```go
span, isRecording := middleware.StartSpan(ctx, "operation_name")
if isRecording {
    // Add attributes to the span
    span.SetAttributes(attribute.String("key", "value"))
}
defer middleware.EndSpan(span, err)
```

### Context Propagation

Always pass the context through your application to maintain the trace context:

```go
func SomeFunction(ctx context.Context, param string) error {
    // The context contains the trace information
    result, err := repository.Query(ctx, param)
    return err
}
```

## Viewing Traces

Depending on the configured exporter, traces can be viewed in different ways:

- **stdout**: Traces are printed to the application logs
- **otlp**: Traces are sent to an OpenTelemetry collector and can be viewed in tools like Jaeger, Zipkin, or other compatible UIs

## Best Practices

1. **Use meaningful span names**: Choose descriptive names that indicate what operation is being performed
2. **Add relevant attributes**: Include information that helps understand the context of the operation
3. **Propagate context**: Always pass the context through your application
4. **Handle errors**: Record errors in spans to make debugging easier
5. **Don't trace everything**: Focus on critical paths and high-value operations
