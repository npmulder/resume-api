#!/bin/bash

# Test script for OpenTelemetry metrics endpoint
echo "Testing OpenTelemetry metrics endpoint..."

# Start the API server in the background
echo "Starting API server..."
go run cmd/api/main.go &
SERVER_PID=$!

# Wait for the server to start
sleep 5

# Make some requests to generate metrics
echo "Making requests to generate metrics..."
curl -s http://localhost:8080/health > /dev/null
curl -s http://localhost:8080/api/v1/profile > /dev/null
curl -s http://localhost:8080/api/v1/experiences > /dev/null
curl -s http://localhost:8080/api/v1/skills > /dev/null

# Wait for metrics to be collected
sleep 2

# Check the metrics endpoint
echo "Checking OpenTelemetry metrics endpoint..."
METRICS=$(curl -s http://localhost:8080/metrics)

# Check if metrics are present (exposed via Prometheus exporter)
if echo "$METRICS" | grep -q "http_requests_total"; then
    echo "✅ HTTP request metrics found"
else
    echo "❌ HTTP request metrics not found"
fi

if echo "$METRICS" | grep -q "http_request_duration_seconds"; then
    echo "✅ HTTP request duration metrics found"
else
    echo "❌ HTTP request duration metrics not found"
fi

if echo "$METRICS" | grep -q "memory_usage_bytes"; then
    echo "✅ Memory usage metrics found"
else
    echo "❌ Memory usage metrics not found"
fi

if echo "$METRICS" | grep -q "goroutines_count"; then
    echo "✅ Goroutines count metrics found"
else
    echo "❌ Goroutines count metrics not found"
fi

# Kill the server
echo "Stopping API server..."
kill $SERVER_PID

echo "OpenTelemetry metrics test completed."
