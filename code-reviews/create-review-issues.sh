#!/bin/bash
# Script to create GitHub issues from code review
# This script provides the commands to create all issues identified in the code review
#
# Usage:
#   ./create-review-issues.sh [review-date-directory]
#
# Example:
#   ./create-review-issues.sh 2025-10-24-comprehensive-review
#
# If no directory is specified, it will use the most recent review directory.

set -e

REPO="huberp/go-vibe"
REVIEW_DIR="${1:-2025-10-24-comprehensive-review}"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "==================================="
echo "Code Review Issues - Creation Script"
echo "==================================="
echo ""
echo "Repository: $REPO"
echo "Review Directory: $REVIEW_DIR"
echo ""
echo "This script will create GitHub issues from the code review findings."
echo "Make sure you have GH_TOKEN set or are authenticated with gh CLI."
echo ""
echo "Note: This script contains a partial implementation covering the first 10 issues."
echo "You may need to extend it for additional issues from CODE_REVIEW_ISSUES.md"
echo ""
read -p "Press Enter to continue or Ctrl+C to cancel..."

# Function to create issue
create_issue() {
    local title="$1"
    local body="$2"
    local labels="$3"
    
    echo "Creating issue: $title"
    echo "$body" | gh issue create \
        --repo "$REPO" \
        --title "$title" \
        --body-file - \
        --label "$labels"
    echo "✓ Issue created"
    echo ""
}

# Issue 1
create_issue \
    "Enhance error handling and logging in authentication handler" \
    "## Description
The \`auth_handler.go\` file has opportunities for improved error handling and logging:
1. Database queries don't log errors before returning generic messages to clients
2. Missing structured logging for authentication attempts (success/failure)
3. No rate limiting context for login attempts

## Goals
- Add structured logging for all authentication events
- Log authentication failures with context (IP, email attempt, timestamp)
- Improve error messages returned to clients while maintaining security
- Add request context propagation

## Files to modify
- \`internal/handlers/auth_handler.go\`
- Add tests in \`internal/handlers/auth_handler_test.go\` (new file)

## Acceptance Criteria
- [ ] All database operations log errors with context
- [ ] Authentication attempts are logged (without passwords)
- [ ] Tests validate error scenarios
- [ ] No sensitive data in logs

**Estimated effort**: Medium
**Priority**: High" \
    "enhancement,security,testing,priority-high"

# Issue 2
create_issue \
    "Add comprehensive tests for auth and user handlers" \
    "## Description
Current test coverage for handlers is 68.9%. Missing tests include:
1. No tests for \`auth_handler.go\` at all
2. Edge cases in user handler (concurrent updates, database failures)
3. Authorization logic for owner vs admin access

## Goals
- Create comprehensive test suite for auth handler
- Increase handler test coverage to >90%
- Test all error paths and edge cases
- Add integration tests for authorization logic

## Files to modify
- Create \`internal/handlers/auth_handler_test.go\`
- Enhance \`internal/handlers/user_handler_test.go\`

## Acceptance Criteria
- [ ] Auth handler has >90% test coverage
- [ ] User handler has >90% test coverage
- [ ] All error paths are tested
- [ ] Authorization edge cases are covered

**Estimated effort**: Large
**Priority**: High" \
    "testing,priority-high"

# Issue 3
create_issue \
    "Improve test coverage for middleware components" \
    "## Description
Middleware test coverage is at 49.4%, with gaps in:
1. Logging middleware edge cases (invalid trace context, long paths)
2. Metrics middleware error scenarios
3. Rate limiting edge cases (burst capacity, multiple IPs)

## Goals
- Increase middleware test coverage to >85%
- Test all middleware error paths
- Add tests for middleware interaction/ordering
- Test edge cases and boundary conditions

## Files to modify
- Enhance \`internal/middleware/logging_test.go\` (new file)
- Enhance \`internal/middleware/metrics_test.go\` (new file)
- Enhance existing middleware tests

## Acceptance Criteria
- [ ] Middleware coverage >85%
- [ ] All error paths tested
- [ ] Edge cases covered
- [ ] Middleware interaction tests added

**Estimated effort**: Medium
**Priority**: High" \
    "testing,middleware,priority-high"

# Issue 4
create_issue \
    "Implement owner access validation in user endpoints" \
    "## Description
The \`GetUserByID\` and \`UpdateUser\` endpoints are documented as \"Owner or Admin\" but don't validate that regular users can only access their own data:
1. Missing validation that user_id from JWT matches the resource ID
2. Admin users can access any user, but regular users should only access their own data
3. No tests validating this authorization logic

## Goals
- Add authorization middleware or handler logic to validate owner access
- Ensure regular users can only read/update their own profile
- Ensure admins can access any user
- Add comprehensive tests for authorization

## Files to modify
- \`internal/handlers/user_handler.go\` (add authorization checks)
- \`internal/handlers/user_handler_test.go\` (add authorization tests)
- Optionally create \`internal/middleware/authorization.go\` for reusable logic

## Acceptance Criteria
- [ ] Regular users can only access their own profile
- [ ] Admins can access any profile
- [ ] Unauthorized access returns 403 Forbidden
- [ ] Tests validate all authorization scenarios

**Estimated effort**: Medium
**Priority**: High" \
    "security,enhancement,priority-high"

# Issue 5
create_issue \
    "Refactor duplicate route definitions for API versioning" \
    "## Description
In \`routes.go\`, there's significant code duplication between v1 routes and legacy routes (lines 99-140):
1. Same handlers registered twice
2. Same middleware applied twice
3. Maintenance burden when updating routes

## Goals
- Create a helper function to register routes once
- Remove duplication while maintaining backward compatibility
- Ensure tests cover both versioned and legacy routes

## Files to modify
- \`internal/routes/routes.go\`
- \`internal/routes/routes_test.go\`

## Acceptance Criteria
- [ ] Route registration logic is DRY (Don't Repeat Yourself)
- [ ] Both v1 and legacy routes still work
- [ ] Tests validate both route versions
- [ ] Code is more maintainable

**Estimated effort**: Small
**Priority**: Medium" \
    "refactoring,code-quality"

# Issue 6
create_issue \
    "Implement JWT token refresh endpoint and logic" \
    "## Description
Current JWT implementation has limitations:
1. Tokens expire after 24 hours with no refresh mechanism
2. Users must re-authenticate to get a new token
3. No \"remember me\" or refresh token functionality

## Goals
- Implement refresh token mechanism (optional refresh tokens)
- Add \`/v1/refresh\` endpoint
- Store refresh tokens securely (in database with expiration)
- Add tests for token refresh logic

## Files to modify
- Create \`internal/handlers/refresh_handler.go\`
- Add \`RefreshToken\` model to \`internal/models/\`
- Update \`pkg/utils/auth.go\` with refresh token generation
- Update \`internal/routes/routes.go\`
- Add tests

## Acceptance Criteria
- [ ] Refresh token endpoint implemented
- [ ] Refresh tokens stored securely
- [ ] Old refresh tokens invalidated on use
- [ ] Comprehensive tests for refresh logic
- [ ] Documentation updated

**Estimated effort**: Large
**Priority**: Medium" \
    "enhancement,feature"

# Issue 7
create_issue \
    "Enhance configuration loading with better error reporting" \
    "## Description
The configuration package has areas for improvement:
1. \`log.Fatalf\` in library code (\`config.go:99\`) violates best practices
2. Missing validation for required configuration values
3. No graceful handling of malformed YAML
4. Limited error context when configuration fails

## Goals
- Return errors instead of using \`log.Fatalf\`
- Validate required configuration fields
- Provide helpful error messages for configuration issues
- Add configuration validation tests

## Files to modify
- \`pkg/config/config.go\`
- \`pkg/config/config_test.go\`
- \`cmd/server/main.go\` (handle config errors)

## Acceptance Criteria
- [ ] No \`log.Fatalf\` in library code
- [ ] Required fields are validated
- [ ] Clear error messages for config issues
- [ ] Tests validate error cases
- [ ] Backward compatible with existing configs

**Estimated effort**: Medium
**Priority**: Medium" \
    "refactoring,code-quality,configuration"

# Issue 8
create_issue \
    "Expose database connection pool metrics to Prometheus" \
    "## Description
While database health checks exist, connection pool metrics aren't exposed:
1. No visibility into connection pool exhaustion
2. Can't monitor connection wait times
3. Missing metrics: active connections, idle connections, wait count, wait duration

## Goals
- Add Prometheus metrics for database connection pool stats
- Expose metrics like: open connections, in use, idle, wait count
- Add Grafana dashboard example in documentation
- Add tests for metrics collection

## Files to modify
- \`internal/middleware/metrics.go\`
- \`docs/observability/METRICS.md\`
- Add example Grafana dashboard

## Acceptance Criteria
- [ ] Connection pool metrics exposed
- [ ] Metrics follow Prometheus naming conventions
- [ ] Documentation includes metric descriptions
- [ ] Example Grafana dashboard provided

**Estimated effort**: Medium
**Priority**: Medium" \
    "observability,enhancement"

# Issue 9
create_issue \
    "Add graceful shutdown for server and database connections" \
    "## Description
The main server doesn't implement graceful shutdown:
1. Server stops immediately when receiving SIGTERM/SIGINT
2. In-flight requests may be terminated abruptly
3. Database connections aren't closed cleanly
4. No timeout for shutdown operations

## Goals
- Implement graceful shutdown with signal handling
- Allow in-flight requests to complete (with timeout)
- Close database connections cleanly
- Add shutdown timeout configuration
- Add tests for shutdown behavior

## Files to modify
- \`cmd/server/main.go\`
- Add shutdown tests

## Acceptance Criteria
- [ ] Server handles SIGTERM and SIGINT gracefully
- [ ] In-flight requests complete before shutdown
- [ ] Database connections closed cleanly
- [ ] Configurable shutdown timeout
- [ ] Logs shutdown progress

**Estimated effort**: Medium
**Priority**: Medium" \
    "enhancement,reliability"

# Issue 10
create_issue \
    "Create comprehensive API validation rules documentation" \
    "## Description
While Swagger docs exist, validation rules aren't fully documented:
1. No clear documentation of validation rules (min/max lengths, formats)
2. Error responses for validation failures aren't documented
3. Examples don't show validation error responses

## Goals
- Document all validation rules in Swagger annotations
- Add examples of validation error responses
- Create API validation guide in docs
- Ensure Swagger UI shows validation constraints

## Files to modify
- Update Swagger annotations in handlers
- Create \`docs/api/VALIDATION.md\`
- Regenerate Swagger documentation

## Acceptance Criteria
- [ ] All validation rules documented in Swagger
- [ ] Validation error examples in Swagger
- [ ] Validation guide created
- [ ] Swagger UI displays constraints

**Estimated effort**: Small
**Priority**: Low" \
    "documentation"

# Continue with remaining issues...
# Issues 11-20 follow same pattern
# 
# To see all 20 issues with complete details, refer to:
# $SCRIPT_DIR/$REVIEW_DIR/CODE_REVIEW_ISSUES.md

echo "==================================="
echo "✓ First 10 issues created successfully!"
echo "==================================="
echo ""
echo "Note: This script contains a partial implementation."
echo "For issues 11-20, please refer to:"
echo "  $SCRIPT_DIR/$REVIEW_DIR/CODE_REVIEW_ISSUES.md"
echo ""
echo "You can extend this script by following the same pattern above."
echo ""
echo "View created issues at: https://github.com/$REPO/issues"
echo ""
echo "For help creating issues manually, see:"
echo "  $SCRIPT_DIR/$REVIEW_DIR/QUICK_START_ISSUES.md"
