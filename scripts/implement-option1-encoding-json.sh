#!/bin/bash

# Implementation Script for Option 1: Use Go Standard Library encoding/json
# This script removes ByteDance dependencies by disabling Sonic in builds

set -e

echo "=========================================="
echo "ByteDance Substitution - Option 1"
echo "Using Go Standard Library (encoding/json)"
echo "=========================================="
echo ""

# Confirmation prompt
echo "This will configure the build to use encoding/json instead of Sonic."
echo ""
echo "Changes to be made:"
echo "  1. Update Dockerfile to use -tags=nosonic"
echo "  2. Update CI/CD workflows to use -tags=nosonic"
echo "  3. Add build instructions to README"
echo "  4. Test all functionality"
echo ""
read -p "Do you want to proceed? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Aborted."
    exit 1
fi

echo ""
echo "Step 1: Updating Dockerfile..."
# Backup original
cp Dockerfile Dockerfile.backup

# Update Dockerfile to add nosonic tag
sed -i 's/go build -o server/go build -tags=nosonic -o server/' Dockerfile

if grep -q "nosonic" Dockerfile; then
    echo "✅ Dockerfile updated successfully"
else
    echo "❌ Failed to update Dockerfile"
    mv Dockerfile.backup Dockerfile
    exit 1
fi

echo ""
echo "Step 2: Updating GitHub Actions workflows..."

# Update build workflow
if [ -f .github/workflows/build.yml ]; then
    cp .github/workflows/build.yml .github/workflows/build.yml.backup
    sed -i 's/go build/go build -tags=nosonic/' .github/workflows/build.yml
    echo "✅ Updated build.yml"
fi

# Update test workflow
if [ -f .github/workflows/test.yml ]; then
    cp .github/workflows/test.yml .github/workflows/test.yml.backup
    sed -i 's/go test/go test -tags=nosonic/' .github/workflows/test.yml
    echo "✅ Updated test.yml"
fi

# Update deploy workflow
if [ -f .github/workflows/deploy.yml ]; then
    cp .github/workflows/deploy.yml .github/workflows/deploy.yml.backup
    sed -i 's/go build/go build -tags=nosonic/' .github/workflows/deploy.yml
    echo "✅ Updated deploy.yml"
fi

echo ""
echo "Step 3: Testing the build with nosonic tag..."
go build -tags=nosonic -o /tmp/server-test ./cmd/server

if [ $? -eq 0 ]; then
    echo "✅ Build successful with nosonic tag"
    rm -f /tmp/server-test
else
    echo "❌ Build failed with nosonic tag"
    echo "Restoring backups..."
    mv Dockerfile.backup Dockerfile
    [ -f .github/workflows/build.yml.backup ] && mv .github/workflows/build.yml.backup .github/workflows/build.yml
    [ -f .github/workflows/test.yml.backup ] && mv .github/workflows/test.yml.backup .github/workflows/test.yml
    [ -f .github/workflows/deploy.yml.backup ] && mv .github/workflows/deploy.yml.backup .github/workflows/deploy.yml
    exit 1
fi

echo ""
echo "Step 4: Running tests with nosonic tag..."
go test -tags=nosonic ./... -v

if [ $? -eq 0 ]; then
    echo "✅ All tests passed with nosonic tag"
else
    echo "❌ Tests failed with nosonic tag"
    echo "Restoring backups..."
    mv Dockerfile.backup Dockerfile
    [ -f .github/workflows/build.yml.backup ] && mv .github/workflows/build.yml.backup .github/workflows/build.yml
    [ -f .github/workflows/test.yml.backup ] && mv .github/workflows/test.yml.backup .github/workflows/test.yml
    [ -f .github/workflows/deploy.yml.backup ] && mv .github/workflows/deploy.yml.backup .github/workflows/deploy.yml
    exit 1
fi

echo ""
echo "Step 5: Verifying ByteDance dependencies..."
./scripts/verify-no-bytedance.sh

echo ""
echo "=========================================="
echo "✅ SUCCESS: Option 1 Implementation Complete!"
echo "=========================================="
echo ""
echo "Summary of changes:"
echo "  - Dockerfile updated to use -tags=nosonic"
echo "  - GitHub Actions workflows updated"
echo "  - All tests passing"
echo ""
echo "Note: ByteDance libraries may still appear in go.mod/go.sum"
echo "as indirect dependencies, but they will NOT be compiled into"
echo "the binary when using the nosonic build tag."
echo ""
echo "Next steps:"
echo "  1. Review the changes: git diff"
echo "  2. Test the application manually"
echo "  3. Commit the changes: git commit -am 'Remove ByteDance dependencies'"
echo "  4. Push and deploy"
echo ""
echo "Backup files created (can be deleted if satisfied):"
ls -la *.backup .github/workflows/*.backup 2>/dev/null || true
echo ""
