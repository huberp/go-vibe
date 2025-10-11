#!/bin/bash

# Local Kubernetes Cluster Cleanup Script
# This script cleans up the local Kind cluster

set -e

echo "==================================="
echo "Local Kubernetes Cluster Cleanup"
echo "==================================="
echo ""

CLUSTER_NAME="go-vibe-local"
NAMESPACE="production"
APP_NAME="myapp"

# Check if cluster exists
if ! kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    echo "Cluster '${CLUSTER_NAME}' does not exist. Nothing to clean up."
    exit 0
fi

echo "What would you like to clean up?"
echo "  1) Only the application deployment (keeps cluster and database)"
echo "  2) Entire namespace (removes app and database)"
echo "  3) Entire cluster (complete cleanup)"
echo ""
read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        echo "Uninstalling application..."
        kubectl config use-context "kind-${CLUSTER_NAME}" 2>/dev/null || true
        helm uninstall "${APP_NAME}" -n "${NAMESPACE}" 2>/dev/null || echo "App not found or already removed"
        echo "Application uninstalled."
        ;;
    2)
        echo "Deleting namespace '${NAMESPACE}'..."
        kubectl config use-context "kind-${CLUSTER_NAME}" 2>/dev/null || true
        kubectl delete namespace "${NAMESPACE}" --ignore-not-found=true
        echo "Namespace deleted."
        ;;
    3)
        echo "Deleting entire cluster '${CLUSTER_NAME}'..."
        kind delete cluster --name "${CLUSTER_NAME}"
        echo "Cluster deleted."
        ;;
    *)
        echo "Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "==================================="
echo "Cleanup Complete!"
echo "==================================="
