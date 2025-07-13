# Resume API

A Go-based REST API that serves resume/CV data from PostgreSQL.

## Features

- RESTful API for resume data
- PostgreSQL database for storage
- Containerized with Docker
- Kubernetes deployment with Helm
- External access via Gateway API (HttpRoute)

## Development

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker (optional, for containerized database)

### Getting Started

1. Clone the repository
2. Copy the example environment file: `cp .env.example .env`
3. Edit the `.env` file with your specific settings
4. Start the database: `make dev-db`
5. Apply migrations: `make migrate-up`
6. Run the API: `go run cmd/api/main.go`

For more detailed development instructions, see [docs/development.md](docs/development.md).

## Deployment

### Docker

Build and run the Docker image:

```bash
docker build -t resume-api .
docker run -p 8080:8080 resume-api
```

### Kubernetes with Helm

The application can be deployed to Kubernetes using the provided Helm chart:

```bash
# Install the chart
helm install resume-api ./deployments/helm/resume-api

# Upgrade an existing installation
helm upgrade resume-api ./deployments/helm/resume-api

# Uninstall the chart
helm uninstall resume-api
```

#### External Access with HttpRoute

The Helm chart includes support for external access using the Gateway API. This requires the Gateway API CRDs to be installed on your cluster.

The chart creates:
1. A Gateway resource that defines the entry point to the cluster
2. An HttpRoute resource that routes traffic to the resume-api service

To configure the Gateway and HttpRoute, edit the values in `values.yaml`:

```yaml
gateway:
  enabled: true
  name: resume-api-gateway
  namespace: default
  className: istio
  listeners:
    - name: http
      port: 80
      protocol: HTTP
      allowedRoutes:
        namespaces:
          from: Same

httpRoute:
  enabled: true
  name: resume-api-route
  hostnames:
    - "resume-api.example.com"
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
        - name: resume-api
          port: 8080
```

For more information about the Helm chart, see [deployments/helm/resume-api/README.md](deployments/helm/resume-api/README.md).

## Testing

Run the tests:

```bash
# Run all tests
make test

# Run tests without database dependencies
make test-short
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.