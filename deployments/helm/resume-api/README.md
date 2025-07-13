# Resume API Helm Chart

This Helm chart deploys the Resume API application on a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- Gateway API CRDs installed on the cluster (for external access)

## Installing the Chart

To install the chart with the release name `resume-api`:

```bash
helm install resume-api ./deployments/helm/resume-api
```

## Configuration

The following table lists the configurable parameters of the Resume API chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `1` |
| `image.repository` | Image repository | `resume-api` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `8080` |
| `gateway.enabled` | Enable Gateway API for external access | `true` |
| `gateway.name` | Gateway name | `resume-api-gateway` |
| `gateway.namespace` | Gateway namespace | `default` |
| `gateway.className` | Gateway class name | `istio` |
| `httpRoute.enabled` | Enable HttpRoute for external access | `true` |
| `httpRoute.name` | HttpRoute name | `resume-api-route` |
| `httpRoute.hostnames` | HttpRoute hostnames | `["resume-api.example.com"]` |
| `config.environment` | Application environment | `production` |
| `config.server.*` | Server configuration | See `values.yaml` |
| `config.logging.*` | Logging configuration | See `values.yaml` |
| `database.external.enabled` | Use external database | `true` |
| `database.external.*` | External database configuration | See `values.yaml` |
| `database.embedded.enabled` | Deploy embedded PostgreSQL | `false` |
| `database.embedded.*` | Embedded PostgreSQL configuration | See `values.yaml` |
| `resources` | Pod resource requests and limits | See `values.yaml` |
| `securityContext` | Pod security context | See `values.yaml` |
| `nodeSelector` | Node selector | `{}` |
| `tolerations` | Tolerations | `[]` |
| `affinity` | Affinity | `{}` |

## External Database

By default, the chart is configured to use an external database. You can provide the database credentials using an existing secret:

```yaml
database:
  external:
    enabled: true
    existingSecret: "my-db-credentials"
    existingSecretPasswordKey: "password"
```

If `existingSecret` is not provided, a random password will be generated.

## Embedded Database

For development or testing, you can enable the embedded PostgreSQL database:

```yaml
database:
  external:
    enabled: false
  embedded:
    enabled: true
```

## External Access

The chart uses the Kubernetes Gateway API for external access. To enable it:

1. Make sure the Gateway API CRDs are installed on your cluster
2. Configure the Gateway and HttpRoute in `values.yaml`

Example:

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

## Customizing the Chart

To customize the chart, create a `values.yaml` file with your changes and use it when installing the chart:

```bash
helm install resume-api ./deployments/helm/resume-api -f my-values.yaml
```