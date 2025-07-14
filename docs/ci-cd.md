# CI/CD Pipeline Documentation

This document describes the Continuous Integration and Continuous Deployment (CI/CD) pipeline for the Resume API project.

## Overview

The CI/CD pipeline is implemented using GitHub Actions and automates the following processes:
- Code linting and static analysis
- Unit testing
- Integration testing with database
- Building the application
- Building Docker images
- (Optional) Deploying to environments

## Workflow Configuration

The pipeline is defined in `.github/workflows/ci.yml` and is triggered on:
- Pushes to the `main` or `master` branch
- Pull requests targeting the `main` or `master` branch

## Pipeline Jobs

### 1. Lint

This job performs static code analysis to ensure code quality and consistency:
- Sets up Go 1.24
- Installs project dependencies
- Installs golangci-lint
- Runs the linter using `make lint`

### 2. Unit Tests

This job runs tests that don't require external dependencies like databases:
- Sets up Go 1.24
- Installs project dependencies
- Runs short tests using `make test-short`

### 3. Integration Tests

This job tests the application's interaction with external dependencies:
- Sets up Go 1.24
- Starts a PostgreSQL 17 container
- Runs database migrations
- Executes repository tests
- Executes integration tests

### 4. Build

This job builds the application and Docker image:
- Sets up Go 1.24
- Installs project dependencies
- Builds the application using `make build`
- Builds a Docker image using `make docker-build`
- (Optional) Pushes the Docker image to a registry

## Environment Variables

The integration tests use the following environment variables:
- `TEST_DB_HOST`: Database host (localhost)
- `TEST_DB_PORT`: Database port (5433)
- `TEST_DB_NAME`: Database name (resume_api_test)
- `TEST_DB_USER`: Database user (dev)
- `TEST_DB_PASSWORD`: Database password (devpass)

## Extending the Pipeline

### Adding Deployment

To add deployment to the pipeline:

1. Uncomment the Docker registry login and push steps in the build job
2. Configure the appropriate secrets in your GitHub repository:
   - `DOCKER_USERNAME`: Your Docker registry username
   - `DOCKER_PASSWORD`: Your Docker registry password/token

3. Add a deployment job that depends on the build job:

```yaml
deploy:
  name: Deploy
  runs-on: ubuntu-latest
  needs: [build]
  if: github.ref == 'refs/heads/main' || github.ref == 'refs/heads/master'
  steps:
    # Add deployment steps here
    # For example, using kubectl to deploy to Kubernetes
```

### Adding Code Coverage

To add code coverage reporting:

1. Modify the test jobs to generate coverage reports
2. Add a step to upload the coverage reports to a service like Codecov

```yaml
- name: Run tests with coverage
  run: go test ./... -coverprofile=coverage.txt -covermode=atomic

- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
  with:
    file: ./coverage.txt
```

## Troubleshooting

### Common Issues

1. **Tests failing in CI but passing locally**
   - Check for environment differences
   - Ensure all required services are properly configured
   - Check for race conditions or timing issues

2. **Docker build failing**
   - Verify the Dockerfile is correct
   - Check if all required files are included in the build context
   - Ensure Docker is properly installed in the CI environment

3. **Linting errors**
   - Run `make lint` locally to identify and fix issues
   - Consider adding a pre-commit hook to catch linting errors before pushing

## Best Practices

1. **Keep the pipeline fast**
   - Use caching for dependencies and build artifacts
   - Run jobs in parallel when possible
   - Use short tests for quick feedback

2. **Secure sensitive information**
   - Use GitHub Secrets for sensitive values
   - Never hardcode credentials in workflow files

3. **Monitor pipeline performance**
   - Regularly review workflow execution times
   - Optimize slow steps
   - Consider using GitHub Actions workflow visualizations