#!/bin/bash

# Local Kubernetes Cluster Setup Script
# This script sets up a local Kind cluster for testing the deployment

set -e

echo "==================================="
echo "Local Kubernetes Cluster Setup"
echo "==================================="
echo ""

# Check if Kind is installed
if ! command -v kind &> /dev/null; then
    echo "Error: Kind is not installed. Please install Kind first."
    echo "Visit: https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo "Error: Docker is not running. Please start Docker first."
    exit 1
fi

# Cluster name
CLUSTER_NAME="go-vibe-local"

# Check if cluster already exists
if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    echo "Cluster '${CLUSTER_NAME}' already exists."
    read -p "Do you want to delete and recreate it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Deleting existing cluster..."
        kind delete cluster --name "${CLUSTER_NAME}"
    else
        echo "Using existing cluster."
        kubectl cluster-info --context "kind-${CLUSTER_NAME}"
        exit 0
    fi
fi

# Create Kind cluster with custom configuration
echo "Creating Kind cluster '${CLUSTER_NAME}'..."
cat <<EOF | kind create cluster --name "${CLUSTER_NAME}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 30080
    protocol: TCP
  - containerPort: 30543
    hostPort: 30543
    protocol: TCP
EOF

# Wait for cluster to be ready
echo "Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=60s

# Verify cluster
echo ""
echo "Cluster created successfully!"
kubectl cluster-info --context "kind-${CLUSTER_NAME}"

echo ""
echo "==================================="
echo "Cluster setup complete!"
echo "==================================="
echo ""
echo "Cluster name: ${CLUSTER_NAME}"
echo "Context: kind-${CLUSTER_NAME}"
echo ""
echo "Next steps:"
echo "  1. Run './scripts/local-k8s-deploy.sh' to deploy the application"
echo "  2. Use 'kubectl get pods -n production' to check pod status"
echo "  3. Use 'kubectl logs -n production <pod-name>' to view logs"
echo ""
