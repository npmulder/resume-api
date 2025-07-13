#!/bin/bash
set -e

echo "Installing swag CLI tool..."
go install github.com/swaggo/swag/cmd/swag@latest

echo "Generating Swagger documentation..."
swag init -g cmd/api/main.go -o docs

echo "Swagger documentation generated successfully!"
echo "You can access the Swagger UI at http://localhost:8080/swagger/index.html when the API is running."