# Local Kubernetes Deployment Troubleshooting

This guide helps troubleshoot common issues when deploying to a local Kind cluster.

## Prerequisites

Before running the local deployment scripts, ensure you have:

1. **Docker** installed and running
   ```bash
   docker info
   ```

2. **Kind** (Kubernetes in Docker) installed
   ```bash
   kind version
   ```

3. **kubectl** installed
   ```bash
   kubectl version --client
   ```

4. **Helm** installed
   ```bash
   helm version
   ```

## Common Issues

### 1. DNS Resolution Issues

**Symptom**: Pods cannot resolve service names (e.g., `postgres`, `kubernetes`)

**Cause**: CoreDNS not working properly due to network restrictions

**Solutions**:

a) **Check CoreDNS status**:
```bash
kubectl get pods -n kube-system | grep coredns
kubectl logs -n kube-system <coredns-pod-name>
```

b) **Restart CoreDNS**:
```bash
kubectl rollout restart deployment/coredns -n kube-system
```

c) **Use IP addresses instead of DNS** (workaround):
   - Get the service IP: `kubectl get svc postgres -n production`
   - Update database URL to use IP instead of hostname

### 2. Image Pull Issues

**Symptom**: Pods stuck in `ImagePullBackOff` or `ErrImagePull`

**Cause**: Network restrictions preventing image pulls from within the Kind cluster

**Solution**: Pre-load images before deployment

```bash
# Pull images locally
docker pull postgres:13-alpine
docker pull <your-app-image>

# Load into Kind cluster
kind load docker-image postgres:13-alpine --name go-vibe-local
kind load docker-image <your-app-image> --name go-vibe-local
```

The deployment scripts handle this automatically for postgres and the application.

### 3. ServiceMonitor CRD Missing

**Symptom**: `no matches for kind "ServiceMonitor"` error

**Cause**: Prometheus Operator CRDs not installed in the cluster

**Solution**: Disable ServiceMonitor in deployment

```bash
helm install myapp ./helm/myapp --set serviceMonitor.enabled=false ...
```

The deployment scripts include this flag automatically.

### 4. Application CrashLoopBackOff

**Symptom**: Application pod continuously restarting

**Diagnosis**:
```bash
kubectl logs -n production <pod-name>
kubectl describe pod -n production <pod-name>
```

**Common Causes**:
- Database connection failure (DNS issue)
- Missing environment variables
- Application errors

**Solutions**:
- Check logs for specific errors
- Verify database is running: `kubectl get pods -n production | grep postgres`
- Check environment variables in deployment

### 5. Helm Deployment Timeout

**Symptom**: `context deadline exceeded` during helm install/upgrade

**Cause**: Pods not becoming ready within timeout period

**Diagnosis**:
```bash
kubectl get pods -n production
kubectl describe pod -n production <pod-name>
```

**Solutions**:
- Increase timeout: `helm install ... --timeout=10m`
- Check pod status and logs
- Fix underlying issues (DNS, image pull, etc.)

## Testing Deployments

### Minimal Test (No Database)

For testing Helm chart syntax and Kubernetes manifests without database dependency:

```bash
# Dry-run to validate manifests
helm install myapp ./helm/myapp \
  --dry-run \
  --debug \
  --namespace production

# Template to see generated YAML
helm template myapp ./helm/myapp \
  --namespace production \
  --set serviceMonitor.enabled=false
```

### Full Local Deployment

Follow these steps in order:

1. **Create cluster**:
   ```bash
   ./scripts/local-k8s-setup.sh
   ```

2. **Pre-pull images** (if network issues):
   ```bash
   docker pull postgres:13-alpine
   ```

3. **Deploy**:
   ```bash
   ./scripts/local-k8s-deploy.sh
   ```

4. **Check status**:
   ```bash
   kubectl get all -n production
   kubectl logs -n production -l app=myapp
   ```

5. **Access application**:
   ```bash
   kubectl port-forward -n production svc/myapp 8080:8080
   curl http://localhost:8080/health
   ```

## Environment-Specific Issues

### GitHub Actions / CI Environments

Some CI environments have network restrictions that prevent:
- DNS resolution within Kind clusters
- Image pulls from Docker Hub within Kind
- TLS certificate verification

**Workarounds**:
1. Pre-build and load all images before deployment
2. Use `imagePullPolicy: Never` for all deployments
3. Consider using `hostNetwork: true` for DNS (security trade-off)
4. Use simpler test approaches (e.g., `helm template` validation)

### Corporate Networks

If behind a corporate proxy or firewall:

1. Configure Docker proxy:
   ```json
   {
     "proxies": {
       "default": {
         "httpProxy": "http://proxy:port",
         "httpsProxy": "http://proxy:port"
       }
     }
   }
   ```

2. Configure Kind with proxy:
   ```yaml
   # kind-config.yaml
   kind: Cluster
   apiVersion: kind.x-k8s.io/v1alpha4
   nodes:
   - role: control-plane
     kubeadmConfigPatches:
     - |
       kind: ClusterConfiguration
       apiServer:
         extraArgs:
           "feature-gates": "AllAlpha=false"
   ```

## Validation Without Full Deployment

If local deployment is not possible due to environment restrictions:

### 1. Lint Helm Chart
```bash
helm lint ./helm/myapp
```

### 2. Validate Kubernetes Manifests
```bash
helm template myapp ./helm/myapp --namespace production | kubectl apply --dry-run=client -f -
```

### 3. Test with Docker Compose
```bash
docker-compose up
./test-api.sh
```

### 4. Manual Deployment Testing
```bash
# Build Docker image
docker build -t myapp:test .

# Run with environment variables
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:password@host:5432/myapp" \
  -e JWT_SECRET="test-secret" \
  myapp:test
```

## Getting Help

If issues persist:

1. Check cluster status:
   ```bash
   kubectl cluster-info
   kubectl get nodes
   kubectl get pods --all-namespaces
   ```

2. Collect logs:
   ```bash
   kubectl logs -n kube-system -l k8s-app=kube-dns
   kubectl logs -n production -l app=myapp
   ```

3. Describe resources:
   ```bash
   kubectl describe pod -n production <pod-name>
   kubectl describe deployment -n production myapp
   ```

4. Export events:
   ```bash
   kubectl get events -n production --sort-by='.lastTimestamp'
   ```

## Cleanup

To remove the local cluster and start fresh:

```bash
# Option 1: Use cleanup script
./scripts/local-k8s-cleanup.sh

# Option 2: Manual cleanup
kind delete cluster --name go-vibe-local
```

## Alternative Testing Approaches

If Kind cluster doesn't work in your environment:

1. **Minikube**: Alternative local Kubernetes
   ```bash
   minikube start
   minikube addons enable metrics-server
   ```

2. **k3d**: Lightweight Kubernetes in Docker
   ```bash
   k3d cluster create go-vibe
   ```

3. **Docker Desktop Kubernetes**: Built-in Kubernetes
   - Enable in Docker Desktop settings
   - Use context: `docker-desktop`

4. **Cloud-based testing**: Use managed Kubernetes
   - Google Kubernetes Engine (GKE)
   - Amazon Elastic Kubernetes Service (EKS)
   - Azure Kubernetes Service (AKS)
   - DigitalOcean Kubernetes

Each has trade-offs in complexity, cost, and environment compatibility.
