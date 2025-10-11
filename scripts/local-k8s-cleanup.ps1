# Local Kubernetes Cluster Cleanup Script (PowerShell)
# This script cleans up the local Kind cluster

$ErrorActionPreference = "Stop"

Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Local Kubernetes Cluster Cleanup" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

$CLUSTER_NAME = "go-vibe-local"
$NAMESPACE = "production"
$APP_NAME = "myapp"

# Check if cluster exists
$existingClusters = kind get clusters 2>$null
if ($existingClusters -notcontains $CLUSTER_NAME) {
    Write-Host "Cluster '$CLUSTER_NAME' does not exist. Nothing to clean up." -ForegroundColor Yellow
    exit 0
}

Write-Host "What would you like to clean up?"
Write-Host "  1) Only the application deployment (keeps cluster and database)"
Write-Host "  2) Entire namespace (removes app and database)"
Write-Host "  3) Entire cluster (complete cleanup)"
Write-Host ""
$choice = Read-Host "Enter your choice (1-3)"

switch ($choice) {
    "1" {
        Write-Host "Uninstalling application..."
        kubectl config use-context "kind-$CLUSTER_NAME" 2>$null
        try {
            helm uninstall $APP_NAME -n $NAMESPACE
        } catch {
            Write-Host "App not found or already removed" -ForegroundColor Yellow
        }
        Write-Host "Application uninstalled." -ForegroundColor Green
    }
    "2" {
        Write-Host "Deleting namespace '$NAMESPACE'..."
        kubectl config use-context "kind-$CLUSTER_NAME" 2>$null
        kubectl delete namespace $NAMESPACE --ignore-not-found=true
        Write-Host "Namespace deleted." -ForegroundColor Green
    }
    "3" {
        Write-Host "Deleting entire cluster '$CLUSTER_NAME'..."
        kind delete cluster --name $CLUSTER_NAME
        Write-Host "Cluster deleted." -ForegroundColor Green
    }
    default {
        Write-Host "Invalid choice. Exiting." -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Cleanup Complete!" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Cyan
