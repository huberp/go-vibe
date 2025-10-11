# Local Kubernetes Cluster Setup Script (PowerShell)
# This script sets up a local Kind cluster for testing the deployment

$ErrorActionPreference = "Stop"

Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Local Kubernetes Cluster Setup" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

# Check if Kind is installed
if (!(Get-Command kind -ErrorAction SilentlyContinue)) {
    Write-Host "Error: Kind is not installed. Please install Kind first." -ForegroundColor Red
    Write-Host "Visit: https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
    exit 1
}

# Check if Docker is running
try {
    docker info | Out-Null
} catch {
    Write-Host "Error: Docker is not running. Please start Docker first." -ForegroundColor Red
    exit 1
}

# Cluster name
$CLUSTER_NAME = "go-vibe-local"

# Check if cluster already exists
$existingClusters = kind get clusters 2>$null
if ($existingClusters -contains $CLUSTER_NAME) {
    Write-Host "Cluster '$CLUSTER_NAME' already exists."
    $response = Read-Host "Do you want to delete and recreate it? (y/N)"
    if ($response -eq 'y' -or $response -eq 'Y') {
        Write-Host "Deleting existing cluster..."
        kind delete cluster --name $CLUSTER_NAME
    } else {
        Write-Host "Using existing cluster."
        kubectl cluster-info --context "kind-$CLUSTER_NAME"
        exit 0
    }
}

# Create Kind cluster with custom configuration
Write-Host "Creating Kind cluster '$CLUSTER_NAME'..."
@"
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
"@ | kind create cluster --name $CLUSTER_NAME --config=-

# Wait for cluster to be ready
Write-Host "Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=60s

# Verify cluster
Write-Host ""
Write-Host "Cluster created successfully!" -ForegroundColor Green
kubectl cluster-info --context "kind-$CLUSTER_NAME"

Write-Host ""
Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Cluster setup complete!" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Cluster name: $CLUSTER_NAME"
Write-Host "Context: kind-$CLUSTER_NAME"
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Run '.\scripts\local-k8s-deploy.ps1' to deploy the application"
Write-Host "  2. Use 'kubectl get pods -n production' to check pod status"
Write-Host "  3. Use 'kubectl logs -n production <pod-name>' to view logs"
Write-Host ""
