# Code Review Issues - Comprehensive List

This document contains all issues identified during the comprehensive code review. Each issue is focused on a particular task that can be addressed by a single coding agent.

## Priority: High

### Issue 1: Improve Error Handling in Auth Handler
**Title**: Enhance error handling and logging in authentication handler

**Description**:
The `auth_handler.go` file has opportunities for improved error handling and logging:
1. Database queries don't log errors before returning generic messages to clients
2. Missing structured logging for authentication attempts (success/failure)
3. No rate limiting context for login attempts

**Goals**:
- Add structured logging for all authentication events
- Log authentication failures with context (IP, email attempt, timestamp)
- Improve error messages returned to clients while maintaining security
- Add request context propagation

**Files to modify**:
- `internal/handlers/auth_handler.go`
- Add tests in `internal/handlers/auth_handler_test.go` (new file)

**Acceptance Criteria**:
- [ ] All database operations log errors with context
- [ ] Authentication attempts are logged (without passwords)
- [ ] Tests validate error scenarios
- [ ] No sensitive data in logs

**Labels**: enhancement, security, testing
**Estimated effort**: Medium

---

### Issue 2: Increase Test Coverage for Handlers
**Title**: Add comprehensive tests for auth and user handlers

**Description**:
Current test coverage for handlers is 68.9%. Missing tests include:
1. No tests for `auth_handler.go` at all
2. Edge cases in user handler (concurrent updates, database failures)
3. Authorization logic for owner vs admin access

**Goals**:
- Create comprehensive test suite for auth handler
- Increase handler test coverage to >90%
- Test all error paths and edge cases
- Add integration tests for authorization logic

**Files to modify**:
- Create `internal/handlers/auth_handler_test.go`
- Enhance `internal/handlers/user_handler_test.go`

**Acceptance Criteria**:
- [ ] Auth handler has >90% test coverage
- [ ] User handler has >90% test coverage
- [ ] All error paths are tested
- [ ] Authorization edge cases are covered

**Labels**: testing, priority-high
**Estimated effort**: Large

---

### Issue 3: Enhance Middleware Test Coverage
**Title**: Improve test coverage for middleware components

**Description**:
Middleware test coverage is at 49.4%, with gaps in:
1. Logging middleware edge cases (invalid trace context, long paths)
2. Metrics middleware error scenarios
3. Rate limiting edge cases (burst capacity, multiple IPs)

**Goals**:
- Increase middleware test coverage to >85%
- Test all middleware error paths
- Add tests for middleware interaction/ordering
- Test edge cases and boundary conditions

**Files to modify**:
- Enhance `internal/middleware/logging_test.go` (new file)
- Enhance `internal/middleware/metrics_test.go` (new file)
- Enhance existing middleware tests

**Acceptance Criteria**:
- [ ] Middleware coverage >85%
- [ ] All error paths tested
- [ ] Edge cases covered
- [ ] Middleware interaction tests added

**Labels**: testing, middleware
**Estimated effort**: Medium

---

### Issue 4: Add Input Validation for Owner Access Control
**Title**: Implement owner access validation in user endpoints

**Description**:
The `GetUserByID` and `UpdateUser` endpoints are documented as "Owner or Admin" but don't validate that regular users can only access their own data:
1. Missing validation that user_id from JWT matches the resource ID
2. Admin users can access any user, but regular users should only access their own data
3. No tests validating this authorization logic

**Goals**:
- Add authorization middleware or handler logic to validate owner access
- Ensure regular users can only read/update their own profile
- Ensure admins can access any user
- Add comprehensive tests for authorization

**Files to modify**:
- `internal/handlers/user_handler.go` (add authorization checks)
- `internal/handlers/user_handler_test.go` (add authorization tests)
- Optionally create `internal/middleware/authorization.go` for reusable logic

**Acceptance Criteria**:
- [ ] Regular users can only access their own profile
- [ ] Admins can access any profile
- [ ] Unauthorized access returns 403 Forbidden
- [ ] Tests validate all authorization scenarios

**Labels**: security, enhancement, priority-high
**Estimated effort**: Medium

---

## Priority: Medium

### Issue 5: Reduce Code Duplication in Routes
**Title**: Refactor duplicate route definitions for API versioning

**Description**:
In `routes.go`, there's significant code duplication between v1 routes and legacy routes (lines 99-140):
1. Same handlers registered twice
2. Same middleware applied twice
3. Maintenance burden when updating routes

**Goals**:
- Create a helper function to register routes once
- Remove duplication while maintaining backward compatibility
- Ensure tests cover both versioned and legacy routes

**Files to modify**:
- `internal/routes/routes.go`
- `internal/routes/routes_test.go`

**Acceptance Criteria**:
- [ ] Route registration logic is DRY (Don't Repeat Yourself)
- [ ] Both v1 and legacy routes still work
- [ ] Tests validate both route versions
- [ ] Code is more maintainable

**Labels**: refactoring, code-quality
**Estimated effort**: Small

---

### Issue 6: Add JWT Token Refresh Mechanism
**Title**: Implement JWT token refresh endpoint and logic

**Description**:
Current JWT implementation has limitations:
1. Tokens expire after 24 hours with no refresh mechanism
2. Users must re-authenticate to get a new token
3. No "remember me" or refresh token functionality

**Goals**:
- Implement refresh token mechanism (optional refresh tokens)
- Add `/v1/refresh` endpoint
- Store refresh tokens securely (in database with expiration)
- Add tests for token refresh logic

**Files to modify**:
- Create `internal/handlers/refresh_handler.go`
- Add `RefreshToken` model to `internal/models/`
- Update `pkg/utils/auth.go` with refresh token generation
- Update `internal/routes/routes.go`
- Add tests

**Acceptance Criteria**:
- [ ] Refresh token endpoint implemented
- [ ] Refresh tokens stored securely
- [ ] Old refresh tokens invalidated on use
- [ ] Comprehensive tests for refresh logic
- [ ] Documentation updated

**Labels**: enhancement, feature
**Estimated effort**: Large

---

### Issue 7: Improve Configuration Error Handling
**Title**: Enhance configuration loading with better error reporting

**Description**:
The configuration package has areas for improvement:
1. `log.Fatalf` in library code (`config.go:99`) violates best practices
2. Missing validation for required configuration values
3. No graceful handling of malformed YAML
4. Limited error context when configuration fails

**Goals**:
- Return errors instead of using `log.Fatalf`
- Validate required configuration fields
- Provide helpful error messages for configuration issues
- Add configuration validation tests

**Files to modify**:
- `pkg/config/config.go`
- `pkg/config/config_test.go`
- `cmd/server/main.go` (handle config errors)

**Acceptance Criteria**:
- [ ] No `log.Fatalf` in library code
- [ ] Required fields are validated
- [ ] Clear error messages for config issues
- [ ] Tests validate error cases
- [ ] Backward compatible with existing configs

**Labels**: refactoring, code-quality, configuration
**Estimated effort**: Medium

---

### Issue 8: Add Database Connection Pool Metrics
**Title**: Expose database connection pool metrics to Prometheus

**Description**:
While database health checks exist, connection pool metrics aren't exposed:
1. No visibility into connection pool exhaustion
2. Can't monitor connection wait times
3. Missing metrics: active connections, idle connections, wait count, wait duration

**Goals**:
- Add Prometheus metrics for database connection pool stats
- Expose metrics like: open connections, in use, idle, wait count
- Add Grafana dashboard example in documentation
- Add tests for metrics collection

**Files to modify**:
- `internal/middleware/metrics.go`
- `docs/observability/METRICS.md`
- Add example Grafana dashboard

**Acceptance Criteria**:
- [ ] Connection pool metrics exposed
- [ ] Metrics follow Prometheus naming conventions
- [ ] Documentation includes metric descriptions
- [ ] Example Grafana dashboard provided

**Labels**: observability, enhancement
**Estimated effort**: Medium

---

### Issue 9: Implement Graceful Shutdown
**Title**: Add graceful shutdown for server and database connections

**Description**:
The main server doesn't implement graceful shutdown:
1. Server stops immediately when receiving SIGTERM/SIGINT
2. In-flight requests may be terminated abruptly
3. Database connections aren't closed cleanly
4. No timeout for shutdown operations

**Goals**:
- Implement graceful shutdown with signal handling
- Allow in-flight requests to complete (with timeout)
- Close database connections cleanly
- Add shutdown timeout configuration
- Add tests for shutdown behavior

**Files to modify**:
- `cmd/server/main.go`
- Add shutdown tests

**Acceptance Criteria**:
- [ ] Server handles SIGTERM and SIGINT gracefully
- [ ] In-flight requests complete before shutdown
- [ ] Database connections closed cleanly
- [ ] Configurable shutdown timeout
- [ ] Logs shutdown progress

**Labels**: enhancement, reliability
**Estimated effort**: Medium

---

### Issue 10: Add API Request/Response Validation Documentation
**Title**: Create comprehensive API validation rules documentation

**Description**:
While Swagger docs exist, validation rules aren't fully documented:
1. No clear documentation of validation rules (min/max lengths, formats)
2. Error responses for validation failures aren't documented
3. Examples don't show validation error responses

**Goals**:
- Document all validation rules in Swagger annotations
- Add examples of validation error responses
- Create API validation guide in docs
- Ensure Swagger UI shows validation constraints

**Files to modify**:
- Update Swagger annotations in handlers
- Create `docs/api/VALIDATION.md`
- Regenerate Swagger documentation

**Acceptance Criteria**:
- [ ] All validation rules documented in Swagger
- [ ] Validation error examples in Swagger
- [ ] Validation guide created
- [ ] Swagger UI displays constraints

**Labels**: documentation
**Estimated effort**: Small

---

## Priority: Low

### Issue 11: Add Request ID to Error Responses
**Title**: Include request ID in error responses for traceability

**Description**:
Error responses don't include request_id:
1. Difficult to correlate client errors with server logs
2. Support teams can't easily trace issues
3. Request ID is generated but not returned to client

**Goals**:
- Add request_id to all error responses
- Update error response structure consistently
- Add tests validating request_id in responses
- Update API documentation

**Files to modify**:
- All handlers (add request_id to error responses)
- Update tests
- Update Swagger documentation

**Acceptance Criteria**:
- [ ] All error responses include request_id
- [ ] Response format is consistent
- [ ] Tests validate request_id presence
- [ ] Documentation updated

**Labels**: enhancement, observability
**Estimated effort**: Medium

---

### Issue 12: Implement Health Check Dependencies
**Title**: Add configurable health check dependencies and readiness logic

**Description**:
Current health checks are simple but could be enhanced:
1. No dependency graph for health checks
2. Can't configure which dependencies are critical vs non-critical
3. Readiness probe treats all dependencies as critical
4. No degraded state support

**Goals**:
- Allow configuring critical vs non-critical dependencies
- Support degraded health state
- Add tests for complex health scenarios
- Document health check configuration

**Files to modify**:
- `pkg/health/provider.go`
- `pkg/health/registry.go`
- `internal/handlers/health_handler.go`
- Add tests

**Acceptance Criteria**:
- [ ] Dependencies can be marked critical/non-critical
- [ ] Degraded state supported
- [ ] Tests cover complex scenarios
- [ ] Documentation updated

**Labels**: enhancement, observability
**Estimated effort**: Large

---

### Issue 13: Add CORS Configuration via Environment Variables
**Title**: Make CORS configuration runtime-configurable

**Description**:
CORS is currently hardcoded in routes.go:
1. AllowOrigins is set to "*" with a comment to configure for production
2. No way to configure CORS without code changes
3. Security risk if deployed to production without changes

**Goals**:
- Add CORS configuration to config files (YAML)
- Support environment variable overrides
- Add validation for CORS configuration
- Document CORS configuration options

**Files to modify**:
- `pkg/config/config.go`
- `internal/routes/routes.go`
- `config/*.yaml` files
- Documentation

**Acceptance Criteria**:
- [ ] CORS configurable via YAML/env vars
- [ ] Default values are secure
- [ ] Configuration validated
- [ ] Documentation updated with examples

**Labels**: configuration, security
**Estimated effort**: Small

---

### Issue 14: Add Structured Logging for Business Events
**Title**: Implement business event logging for audit trails

**Description**:
Current logging focuses on HTTP requests but misses business events:
1. No logging for user creation/deletion
2. No logging for role changes
3. No audit trail for sensitive operations
4. Difficult to track security-relevant events

**Goals**:
- Add structured logging for business events
- Log user CRUD operations with actor information
- Log authentication events (success/failure)
- Add audit log documentation

**Files to modify**:
- `internal/handlers/user_handler.go`
- `internal/handlers/auth_handler.go`
- Create `docs/observability/AUDIT_LOGS.md`

**Acceptance Criteria**:
- [ ] Business events are logged with context
- [ ] Audit trail is complete
- [ ] Logs include actor (who) and action (what)
- [ ] Documentation describes audit log format

**Labels**: enhancement, observability, security
**Estimated effort**: Medium

---

### Issue 15: Optimize Imports and Dependencies
**Title**: Review and optimize Go module dependencies

**Description**:
Current dependencies could be optimized:
1. Some indirect dependencies might be outdated
2. No automated dependency vulnerability scanning
3. go.mod should be regularly updated

**Goals**:
- Run `go mod tidy` and verify dependencies
- Check for security vulnerabilities in dependencies
- Update dependencies to latest stable versions
- Document dependency management practices

**Files to modify**:
- `go.mod`, `go.sum`
- Update `.github/workflows/` to add dependency scanning
- Document dependency update process

**Acceptance Criteria**:
- [ ] All dependencies up to date
- [ ] No known vulnerabilities
- [ ] Dependency scanning in CI
- [ ] Update process documented

**Labels**: maintenance, security
**Estimated effort**: Small

---

### Issue 16: Add Database Migration Testing
**Title**: Implement automated tests for database migrations

**Description**:
Migrations exist but aren't tested:
1. No tests for migration up/down scripts
2. No validation of migration idempotency
3. Risk of migration failures in production

**Goals**:
- Add integration tests for migrations
- Test both up and down migrations
- Validate migration idempotency
- Document migration testing process

**Files to modify**:
- Create `pkg/migration/migration_test.go`
- Add migration test helpers
- Update CI to run migration tests

**Acceptance Criteria**:
- [ ] All migrations have tests
- [ ] Up and down migrations tested
- [ ] Idempotency validated
- [ ] Tests run in CI

**Labels**: testing, database
**Estimated effort**: Medium

---

### Issue 17: Add Dockerfile Security Scanning
**Title**: Integrate container security scanning in CI/CD

**Description**:
Docker images aren't scanned for vulnerabilities:
1. No security scanning in build pipeline
2. Base images might have vulnerabilities
3. No policy for updating base images

**Goals**:
- Add Trivy or Grype security scanning to CI
- Fail builds on high/critical vulnerabilities
- Document vulnerability management process
- Add automated base image updates

**Files to modify**:
- `.github/workflows/build.yml`
- Add security scanning workflow
- Update documentation

**Acceptance Criteria**:
- [ ] Container scanning in CI
- [ ] Builds fail on critical vulnerabilities
- [ ] Process documented
- [ ] Automated updates configured

**Labels**: security, ci-cd, docker
**Estimated effort**: Small

---

### Issue 18: Add Performance Benchmarks
**Title**: Create performance benchmark tests for critical paths

**Description**:
No performance benchmarks exist:
1. Can't track performance regressions
2. No baseline for optimization efforts
3. Critical paths (auth, user CRUD) aren't benchmarked

**Goals**:
- Add Go benchmark tests for critical handlers
- Benchmark authentication flow
- Benchmark database operations
- Document performance expectations

**Files to modify**:
- Create benchmark tests (`*_bench_test.go`)
- Add benchmark running to CI
- Document benchmark results

**Acceptance Criteria**:
- [ ] Benchmarks for critical paths
- [ ] Benchmarks run in CI
- [ ] Baseline performance documented
- [ ] Regression detection in place

**Labels**: testing, performance
**Estimated effort**: Medium

---

### Issue 19: Improve Swagger Documentation Quality
**Title**: Enhance OpenAPI/Swagger documentation completeness

**Description**:
Swagger docs exist but could be improved:
1. Some error responses not documented
2. Request/response examples could be more complete
3. Authentication flow not clearly documented
4. No examples for all endpoints

**Goals**:
- Complete all Swagger annotations
- Add comprehensive examples for all endpoints
- Document all error responses
- Add authentication guide to Swagger UI

**Files to modify**:
- Update Swagger annotations in all handlers
- Regenerate Swagger docs
- Add examples to Swagger

**Acceptance Criteria**:
- [ ] All endpoints fully documented
- [ ] All error codes documented
- [ ] Examples for all endpoints
- [ ] Authentication clearly explained

**Labels**: documentation
**Estimated effort**: Medium

---

### Issue 20: Add Rate Limiting Per User
**Title**: Implement per-user rate limiting for authenticated endpoints

**Description**:
Current rate limiting is per-IP only:
1. Can't limit authenticated users by user ID
2. Users behind same proxy share rate limits
3. No differentiation between user roles

**Goals**:
- Add per-user rate limiting for authenticated requests
- Support different limits for different roles
- Maintain IP-based limits for public endpoints
- Add tests for rate limiting scenarios

**Files to modify**:
- `internal/middleware/ratelimit.go`
- Add user-based rate limiter
- Update configuration for per-user limits
- Add tests

**Acceptance Criteria**:
- [ ] Per-user rate limiting implemented
- [ ] Different limits for different roles
- [ ] IP-based limits still work
- [ ] Configuration documented
- [ ] Tests cover all scenarios

**Labels**: enhancement, security
**Estimated effort**: Large

---

## Summary Statistics

**Total Issues**: 20
- **Priority High**: 4 issues
- **Priority Medium**: 9 issues  
- **Priority Low**: 7 issues

**Effort Distribution**:
- **Small**: 4 issues
- **Medium**: 11 issues
- **Large**: 5 issues

**Categories**:
- Testing: 6 issues
- Security: 5 issues
- Enhancement: 8 issues
- Documentation: 3 issues
- Refactoring: 2 issues
- Configuration: 2 issues
- Observability: 4 issues
- CI/CD: 2 issues

## Instructions for Creating Issues

Use the GitHub web interface or CLI to create issues with the following template:

```
Title: [Copy from above]

Labels: [Copy from above]

Description:
[Copy full description section]

Acceptance Criteria:
[Copy checklist]

Estimated Effort: [Small/Medium/Large]
Priority: [High/Medium/Low]
```

Each issue is scoped to be addressed by a single coding agent and focused on a particular task.
