#!/bin/bash

# Local Kubernetes Deployment Script
# This script deploys the application to a local Kind cluster

set -e

echo "==================================="
echo "Local Kubernetes Deployment"
echo "==================================="
echo ""

# Configuration
CLUSTER_NAME="go-vibe-local"
NAMESPACE="production"
APP_NAME="myapp"
IMAGE_NAME="myapp"
IMAGE_TAG="local"

# Check if Kind cluster exists
if ! kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    echo "Error: Cluster '${CLUSTER_NAME}' does not exist."
    echo "Please run './scripts/local-k8s-setup.sh' first."
    exit 1
fi

# Set kubectl context
echo "Setting kubectl context..."
kubectl config use-context "kind-${CLUSTER_NAME}"

# Build the application locally first
echo ""
echo "Building application binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./server ./cmd/server

# Build Docker image using pre-built binary
echo ""
echo "Building Docker image..."
cat > /tmp/Dockerfile.kind <<'DOCKERFILE'
FROM gcr.io/distroless/static-debian11
COPY server /server
EXPOSE 8080
USER nonroot:nonroot
CMD ["/server"]
DOCKERFILE

docker build -f /tmp/Dockerfile.kind -t "${IMAGE_NAME}:${IMAGE_TAG}" .

# Load image into Kind cluster
echo ""
echo "Loading Docker image into Kind cluster..."
kind load docker-image "${IMAGE_NAME}:${IMAGE_TAG}" --name "${CLUSTER_NAME}"

# Pre-pull and load postgres image to avoid network issues
echo ""
echo "Preparing PostgreSQL image..."
if docker pull postgres:13-alpine --quiet 2>/dev/null; then
    kind load docker-image postgres:13-alpine --name "${CLUSTER_NAME}"
    SKIP_POSTGRES=false
else
    echo "Warning: Could not pull postgres image. Skipping database deployment..."
    SKIP_POSTGRES=true
fi

# Create namespace
echo ""
echo "Creating namespace '${NAMESPACE}'..."
kubectl create namespace "${NAMESPACE}" --dry-run=client -o yaml | kubectl apply -f -

# Deploy PostgreSQL for local testing (if not skipped)
if [ "${SKIP_POSTGRES}" != "true" ]; then
    echo ""
    echo "Deploying PostgreSQL..."
    kubectl apply -n "${NAMESPACE}" -f - <<EOF
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
EOF

    # Wait for PostgreSQL to be ready
    echo ""
    echo "Waiting for PostgreSQL to be ready..."
    kubectl wait --for=condition=available --timeout=120s deployment/postgres -n "${NAMESPACE}" || {
        echo "Warning: PostgreSQL deployment timed out. Checking status..."
        kubectl get pods -n "${NAMESPACE}" | grep postgres
    }
fi

# Deploy application using Helm
echo ""
echo "Deploying application with Helm..."
helm upgrade --install "${APP_NAME}" ./helm/myapp \
  --namespace "${NAMESPACE}" \
  --set image.repository="${IMAGE_NAME}" \
  --set image.tag="${IMAGE_TAG}" \
  --set image.pullPolicy=Never \
  --set database.url="postgres://user:password@postgres:5432/myapp?sslmode=disable" \
  --set jwt.secret="local-test-secret-key" \
  --set autoscaling.enabled=false \
  --set serviceMonitor.enabled=false \
  --set replicaCount=1 \
  --wait \
  --timeout=5m || {
    echo ""
    echo "Warning: Helm deployment encountered issues. Checking status..."
    kubectl get pods -n "${NAMESPACE}"
  }

# Wait for deployment to be ready
echo ""
echo "Waiting for application to be ready..."
kubectl wait --for=condition=available --timeout=120s deployment/${APP_NAME} -n "${NAMESPACE}"

# Get pod status
echo ""
echo "==================================="
echo "Deployment Status"
echo "==================================="
kubectl get pods -n "${NAMESPACE}"

echo ""
echo "==================================="
echo "Services"
echo "==================================="
kubectl get svc -n "${NAMESPACE}"

# Port forward instructions
echo ""
echo "==================================="
echo "Deployment Complete!"
echo "==================================="
echo ""
echo "To access the application, run:"
echo "  kubectl port-forward -n ${NAMESPACE} svc/${APP_NAME} 8080:8080"
echo ""
echo "Then access the API at: http://localhost:8080"
echo ""
echo "Health check: curl http://localhost:8080/health"
echo "Metrics: curl http://localhost:8080/metrics"
echo ""
echo "To view logs:"
echo "  kubectl logs -n ${NAMESPACE} -l app=${APP_NAME} -f"
echo ""
echo "To delete the deployment:"
echo "  helm uninstall ${APP_NAME} -n ${NAMESPACE}"
echo ""

# Cleanup local binary
rm -f ./server /tmp/Dockerfile.kind
