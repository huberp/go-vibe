# Local Kubernetes Deployment Script (PowerShell)
# This script deploys the application to a local Kind cluster

$ErrorActionPreference = "Stop"

Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Local Kubernetes Deployment" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

# Configuration
$CLUSTER_NAME = "go-vibe-local"
$NAMESPACE = "production"
$APP_NAME = "myapp"
$IMAGE_NAME = "myapp"
$IMAGE_TAG = "local"

# Check if Kind cluster exists
$existingClusters = kind get clusters 2>$null
if ($existingClusters -notcontains $CLUSTER_NAME) {
    Write-Host "Error: Cluster '$CLUSTER_NAME' does not exist." -ForegroundColor Red
    Write-Host "Please run '.\scripts\local-k8s-setup.ps1' first."
    exit 1
}

# Set kubectl context
Write-Host "Setting kubectl context..."
kubectl config use-context "kind-$CLUSTER_NAME"

# Build Docker image
Write-Host ""
Write-Host "Building Docker image..."
docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" .

# Load image into Kind cluster
Write-Host ""
Write-Host "Loading Docker image into Kind cluster..."
kind load docker-image "${IMAGE_NAME}:${IMAGE_TAG}" --name $CLUSTER_NAME

# Pre-pull and load postgres image to avoid network issues
Write-Host ""
Write-Host "Preparing PostgreSQL image..."
try {
    docker pull postgres:13-alpine --quiet 2>$null | Out-Null
    kind load docker-image postgres:13-alpine --name $CLUSTER_NAME
    $SKIP_POSTGRES = $false
} catch {
    Write-Host "Warning: Could not pull postgres image. Skipping database deployment..." -ForegroundColor Yellow
    $SKIP_POSTGRES = $true
}

# Create namespace
Write-Host ""
Write-Host "Creating namespace '$NAMESPACE'..."
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Deploy PostgreSQL for local testing (if not skipped)
if (-not $SKIP_POSTGRES) {
    Write-Host ""
    Write-Host "Deploying PostgreSQL..."
    @"
apiVersion: v1
kind: Service
metadata:
  name: postgres
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:13-alpine
        imagePullPolicy: Never
        env:
        - name: POSTGRES_USER
          value: "user"
        - name: POSTGRES_PASSWORD
          value: "password"
        - name: POSTGRES_DB
          value: "myapp"
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        emptyDir: {}
"@ | kubectl apply -n $NAMESPACE -f -

    # Wait for PostgreSQL to be ready
    Write-Host ""
    Write-Host "Waiting for PostgreSQL to be ready..."
    try {
        kubectl wait --for=condition=available --timeout=120s deployment/postgres -n $NAMESPACE
    } catch {
        Write-Host "Warning: PostgreSQL deployment timed out. Checking status..." -ForegroundColor Yellow
        kubectl get pods -n $NAMESPACE | Select-String postgres
    }
}

# Deploy application using Helm
Write-Host ""
Write-Host "Deploying application with Helm..."
helm upgrade --install $APP_NAME ./helm/myapp `
  --namespace $NAMESPACE `
  --set image.repository=$IMAGE_NAME `
  --set image.tag=$IMAGE_TAG `
  --set image.pullPolicy=Never `
  --set database.url="postgres://user:password@postgres:5432/myapp?sslmode=disable" `
  --set jwt.secret="local-test-secret-key" `
  --set autoscaling.enabled=false `
  --set serviceMonitor.enabled=false `
  --set replicaCount=1 `
  --wait `
  --timeout=5m

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "Warning: Helm deployment encountered issues. Checking status..." -ForegroundColor Yellow
    kubectl get pods -n $NAMESPACE
}

# Wait for deployment to be ready
Write-Host ""
Write-Host "Waiting for application to be ready..."
kubectl wait --for=condition=available --timeout=120s deployment/$APP_NAME -n $NAMESPACE

# Get pod status
Write-Host ""
Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Deployment Status" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
kubectl get pods -n $NAMESPACE

Write-Host ""
Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Services" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
kubectl get svc -n $NAMESPACE

# Port forward instructions
Write-Host ""
Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Deployment Complete!" -ForegroundColor Green
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "To access the application, run:"
Write-Host "  kubectl port-forward -n $NAMESPACE svc/$APP_NAME 8080:8080"
Write-Host ""
Write-Host "Then access the API at: http://localhost:8080"
Write-Host ""
Write-Host "Health check: curl http://localhost:8080/health"
Write-Host "Metrics: curl http://localhost:8080/metrics"
Write-Host ""
Write-Host "To view logs:"
Write-Host "  kubectl logs -n $NAMESPACE -l app=$APP_NAME -f"
Write-Host ""
Write-Host "To delete the deployment:"
Write-Host "  helm uninstall $APP_NAME -n $NAMESPACE"
Write-Host ""
