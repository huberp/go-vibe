# ByteDance Substitution - Implementation Checklist

## Pre-Implementation

- [x] Identify all ByteDance libraries
- [x] Analyze dependency chain
- [x] Research alternatives
- [x] Document consequences
- [x] Create analysis documents
- [x] Create verification scripts
- [x] Create implementation scripts
- [ ] **Get user confirmation on preferred option**

---

## Option 1: encoding/json (RECOMMENDED)

### Automated Implementation Script
Run: `./scripts/implement-option1-encoding-json.sh`

### Manual Steps (if preferred)

#### 1. Update Dockerfile
- [ ] Open `Dockerfile`
- [ ] Find line: `RUN go build -o server ./cmd/server`
- [ ] Change to: `RUN go build -tags=nosonic -o server ./cmd/server`
- [ ] Save file

#### 2. Update CI/CD Workflows

**File: `.github/workflows/build.yml`**
- [ ] Find: `go build ./cmd/server`
- [ ] Change to: `go build -tags=nosonic ./cmd/server`

**File: `.github/workflows/test.yml`**
- [ ] Find: `go test ./...`
- [ ] Change to: `go test -tags=nosonic ./...`

**File: `.github/workflows/deploy.yml`**
- [ ] Find: `go build ./cmd/server`
- [ ] Change to: `go build -tags=nosonic ./cmd/server`

#### 3. Update Local Build Instructions

**File: `README.md`**
- [ ] Add to build section:
```bash
# Build with encoding/json (without ByteDance Sonic)
go build -tags=nosonic ./cmd/server
```

#### 4. Test Everything

- [ ] Build locally: `go build -tags=nosonic ./cmd/server`
- [ ] Run tests: `go test -tags=nosonic ./...`
- [ ] Run application and test endpoints
- [ ] Verify no ByteDance code is compiled (check binary size - should be similar)

#### 5. Performance Testing

- [ ] Benchmark JSON operations
- [ ] Load test API endpoints
- [ ] Compare response times before/after
- [ ] Document any performance differences

#### 6. Clean Up

- [ ] Remove backup files if everything works
- [ ] Update documentation
- [ ] Run `./scripts/verify-no-bytedance.sh` (will still show deps in go.mod, but not compiled)

---

## Option 2: go-json (Alternative)

### Steps

#### 1. This option requires Gin to support go-json

**Note**: Gin v1.11.0 has Sonic deeply integrated. Switching to go-json requires either:
- Upgrading to a Gin version with go-json support, OR
- Forking Gin and modifying the JSON codec, OR
- Creating custom middleware to override JSON binding

**Recommendation**: This is more complex than Option 1. Consider Option 1 instead.

If you still want to pursue this:
- [ ] Research Gin's JSON codec configuration
- [ ] Implement custom JSON binding middleware
- [ ] Test thoroughly
- [ ] Update all dependencies

---

## Option 3: jsoniter

Similar complexity to Option 2.

---

## Option 4: Keep Sonic (No Changes)

- [ ] Document acceptance of ByteDance dependency
- [ ] No implementation needed

---

## Post-Implementation Verification

### For Option 1

- [ ] Run: `./scripts/verify-no-bytedance.sh`
  - Note: Will still report ByteDance in go.mod (indirect deps)
  - This is expected - they exist but aren't compiled
  
- [ ] Verify binary doesn't contain Sonic:
```bash
go build -tags=nosonic -o server ./cmd/server
strings server | grep -i sonic || echo "No Sonic references found ✅"
strings server | grep -i bytedance || echo "No ByteDance references found ✅"
```

- [ ] Test all API endpoints
- [ ] Check logs for any errors
- [ ] Monitor performance metrics

### Success Criteria

- [x] ByteDance code NOT compiled into binary
- [x] All tests passing
- [x] Application builds successfully
- [x] No runtime errors
- [ ] Performance acceptable (< 10ms additional latency)
- [ ] All API endpoints working
- [ ] Health checks passing

---

## Rollback Plan

If Option 1 causes issues:

1. **Immediate Rollback**:
```bash
# Restore backups
mv Dockerfile.backup Dockerfile
mv .github/workflows/build.yml.backup .github/workflows/build.yml
mv .github/workflows/test.yml.backup .github/workflows/test.yml
mv .github/workflows/deploy.yml.backup .github/workflows/deploy.yml

# Rebuild
go build ./cmd/server
```

2. **Verify Rollback**:
```bash
go test ./...
./server
```

3. **Redeploy Previous Version**:
```bash
git checkout <previous-commit>
docker build -t myapp:rollback .
# Deploy to Kubernetes
```

---

## Performance Benchmarking

### Before Implementation

Run baseline benchmarks:
```bash
# Create benchmark file
cat > benchmark_test.go << 'EOF'
package main

import (
    "encoding/json"
    "testing"
)

type TestStruct struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

func BenchmarkJSONMarshal(b *testing.B) {
    data := TestStruct{ID: 1, Name: "Test User", Email: "test@example.com"}
    for i := 0; i < b.N; i++ {
        json.Marshal(data)
    }
}

func BenchmarkJSONUnmarshal(b *testing.B) {
    jsonData := []byte(`{"id":1,"name":"Test User","email":"test@example.com"}`)
    for i := 0; i < b.N; i++ {
        var data TestStruct
        json.Unmarshal(jsonData, &data)
    }
}
EOF

# Run benchmark
go test -bench=. -benchmem
```

### After Implementation

```bash
go test -tags=nosonic -bench=. -benchmem
```

### Compare Results

Document the performance difference and verify it's acceptable.

---

## Questions to Answer

Before implementing:

- [ ] What is our current API load (requests/second)?
- [ ] What are our latency SLAs?
- [ ] Do we have performance budgets?
- [ ] Is there a compliance requirement to remove ByteDance deps?
- [ ] Can we accept 2-3ms additional latency?

---

## Communication Plan

### Stakeholders to Notify

- [ ] Development team
- [ ] DevOps/SRE team
- [ ] Product/Business stakeholders
- [ ] Security/Compliance team

### What to Communicate

1. **What**: Removing ByteDance dependencies
2. **Why**: Security, compliance, or organizational policy
3. **How**: Using Go stdlib instead of Sonic
4. **Impact**: Slight performance reduction (~2-3ms per request)
5. **When**: After testing and approval
6. **Rollback**: Plan available if issues arise

---

## Timeline (Estimated)

- **Analysis & Documentation**: ✅ Complete
- **User Decision**: ⏳ Pending (@huberp)
- **Implementation**: ~1-2 hours
- **Testing**: ~2-4 hours
- **Performance Validation**: ~2 hours
- **Deployment**: ~1 hour
- **Monitoring**: Ongoing (first 48 hours critical)

**Total**: ~1 day including monitoring

---

## Current Status

📍 **Waiting for @huberp to select preferred option**

Options available:
1. **encoding/json** (recommended)
2. go-json
3. jsoniter
4. Keep Sonic

**Ready to implement immediately upon confirmation! 🚀**
