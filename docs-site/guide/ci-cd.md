# CI/CD

go-vibe ships with a complete GitHub Actions CI/CD pipeline covering build, test, and deployment — including cross-platform testing that validates the API scripts on both Linux and Windows.

## Workflow Overview

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `build.yml` | Push to `main`/`develop`, PRs | Compile the binary |
| `test.yml` | Push to `main`/`develop`, PRs | Run Go tests with coverage |
| `deploy.yml` | Push to `main`, version tags | Build Docker image, deploy to K8s |
| `scripts-test.yml` | Push, PRs touching `test-api.*` | Test API shell/PowerShell scripts on Linux & Windows |
| `pages.yml` | Push to `main` touching `docs-site/` | Build and deploy this documentation to GitHub Pages |

## Build Workflow (`build.yml`)

```yaml
name: Build

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Build
        run: go build -v ./...

      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: server-binary
          path: ./server
```

## Test Workflow (`test.yml`)

Tests run against a real PostgreSQL database using `ikalnytskyi/action-setup-postgres`:

```yaml
name: Test

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Set up PostgreSQL
        uses: ikalnytskyi/action-setup-postgres@v7
        with:
          username: testuser
          password: testpassword
          database: testdb
          port: 5432
        id: postgres

      - name: Run tests
        env:
          DATABASE_URL: ${{ steps.postgres.outputs.connection-uri }}
          JWT_SECRET: test-secret-key
        run: go test ./... -v -race -coverprofile=coverage.out

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v4
        with:
          file: coverage.out
```

## Windows PostgreSQL Testing

One of the more complex CI challenges was running integration tests on Windows runners where PostgreSQL isn't natively available. This was tracked in [issue #3](https://github.com/huberp/go-vibe/issues/3) and resolved in [issue #4](https://github.com/huberp/go-vibe/issues/4).

The solution uses `ikalnytskyi/action-setup-postgres` which works seamlessly on both `ubuntu-latest` and `windows-latest`:

```yaml
# From scripts-test.yml — cross-platform PostgreSQL setup
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest]

runs-on: ${{ matrix.os }}

steps:
  - name: Set up PostgreSQL
    uses: ikalnytskyi/action-setup-postgres@v7
    with:
      username: testuser
      password: testpassword
      database: testdb
      port: 5432
    id: postgres

  - name: Start server (Linux)
    if: runner.os == 'Linux'
    env:
      DATABASE_URL: ${{ steps.postgres.outputs.connection-uri }}
      JWT_SECRET: test-secret
    run: |
      go run ./cmd/server &
      echo $! > server.pid
      sleep 3

  - name: Start server (Windows)
    if: runner.os == 'Windows'
    env:
      DATABASE_URL: ${{ steps.postgres.outputs.connection-uri }}
      JWT_SECRET: test-secret
    run: |
      $proc = Start-Process -FilePath "go" -ArgumentList "run","./cmd/server" -RedirectStandardOutput "server.log" -RedirectStandardError "server-err.log" -PassThru
      $proc.Id | Out-File server.pid
      Start-Sleep -Seconds 5
    shell: pwsh
```

## PID File Handling (Windows)

Stopping a background server process on Windows required special handling due to differences in PID file format and process management. This was improved in [issue #5](https://github.com/huberp/go-vibe/issues/5) and [issue #6](https://github.com/huberp/go-vibe/issues/6).

```powershell
# Robust server stop on Windows — handles BOM and whitespace in PID files
- name: Stop server (Windows)
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    if (Test-Path server.pid) {
      # Read PID, stripping BOM and whitespace that Out-File can introduce
      $raw = Get-Content server.pid -Raw
      $pid = $raw -replace '[\xEF\xBB\xBF\r\n\s]', ''
      if ($pid -match '^\d+$') {
        Stop-Process -Id ([int]$pid) -Force -ErrorAction SilentlyContinue
        Write-Host "Stopped server PID $pid"
      } else {
        Write-Warning "Could not parse PID: '$pid'"
      }
    }
```

::: tip Cross-Platform Scripts
The repo includes both `test-api.sh` (bash) and `test-api.ps1` (PowerShell) for testing the API. Both are validated in CI on their respective platforms.
:::

## Deploy Workflow (`deploy.yml`)

Triggered on pushes to `main` and version tags:

```yaml
name: Deploy

on:
  push:
    branches: [main]
    tags: ['v*.*.*']

jobs:
  docker:
    runs-on: ubuntu-latest
    permissions:
      packages: write
    steps:
      - uses: actions/checkout@v4

      - name: Log in to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          push: true
          tags: ghcr.io/${{ github.repository }}:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy:
    needs: docker
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Deploy with Helm
        env:
          KUBECONFIG: ${{ secrets.KUBECONFIG }}
        run: |
          helm upgrade --install go-vibe ./helm/myapp \
            --set image.tag=${{ github.sha }} \
            --wait --timeout 5m
```

## GitHub Pages Docs Deployment

The documentation you're reading is automatically deployed when anything in `docs-site/` changes on `main`. See the [`pages.yml` workflow](https://github.com/huberp/go-vibe/blob/main/.github/workflows/pages.yml) for the full definition.

```bash
# Trigger a manual docs rebuild
gh workflow run pages.yml
```

## Branch Strategy

```
main          ← production deployments, protected branch
  └── develop ← integration branch
        └── feature/<issue>-description  ← feature branches
        └── fix/<issue>-description      ← bug fix branches
```

PRs to `main` require:
- All CI checks green (build + test)
- At least one maintainer approval
- No merge conflicts

## Secrets Required

| Secret | Used by | Description |
|--------|---------|-------------|
| `GITHUB_TOKEN` | `deploy.yml` | Auto-provided by GitHub Actions |
| `KUBECONFIG` | `deploy.yml` | Kubernetes cluster credentials |
| `CODECOV_TOKEN` | `test.yml` | Codecov upload token (optional for public repos) |
