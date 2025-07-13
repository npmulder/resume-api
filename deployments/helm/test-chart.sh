#!/bin/bash
set -e

# Script to test the resume-api Helm chart

# Change to the directory containing this script
cd "$(dirname "$0")"

echo "Linting the Helm chart..."
helm lint resume-api

echo "Validating the Helm chart templates..."
helm template resume-api resume-api

echo "Validating the Helm chart with Kubernetes..."
helm template resume-api resume-api | kubectl apply --dry-run=client -f -

echo "All tests passed!"