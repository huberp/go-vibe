# Code Review Instructions

This document provides comprehensive guidelines for conducting in-depth code reviews of the go-vibe project.

## Purpose

Code reviews serve to:
- Identify areas of improvement across the codebase
- Create actionable tasks for incremental enhancement
- Maintain code quality and security standards
- Track technical debt and improvement opportunities
- Guide development priorities

## Review Process

### 1. Preparation

**Before Starting**:
- Ensure you have a clean clone of the repository
- Run all tests to establish baseline: `go test ./... -v`
- Check test coverage: `go test ./... -coverprofile=coverage.out`
- Review recent commits and change history
- Understand the project's current state and goals

### 2. Comprehensive Analysis

Analyze the following categories:

#### A. Architecture & Code Organization
- Package structure and separation of concerns
- Use of design patterns (repository, dependency injection)
- Layer separation (handlers, middleware, models, repository)
- Code organization and modularity

#### B. Security
- Authentication and authorization implementation
- Password hashing and credential management
- Input validation and sanitization
- SQL injection prevention
- CORS configuration
- Secrets management
- Rate limiting and abuse prevention

#### C. Testing
- Test coverage across all packages
- Test quality (unit, integration, edge cases)
- Test organization and maintainability
- Mock usage and test isolation
- TDD compliance

#### D. Error Handling & Logging
- Error propagation and handling
- Structured logging implementation
- Log levels and context
- Error messages (client vs server)
- Request tracing and correlation

#### E. API Design
- RESTful conventions
- HTTP method usage
- Status code correctness
- Request/response validation
- API documentation (Swagger/OpenAPI)
- Versioning strategy

#### F. Database & Persistence
- Repository pattern implementation
- Query optimization
- Connection pooling
- Migration management
- Data model design

#### G. Configuration Management
- Environment-based configuration
- Secrets handling
- Default values
- Validation of required fields
- Configuration documentation

#### H. Observability
- Metrics exposure (Prometheus)
- Health check implementation
- Distributed tracing (OpenTelemetry)
- Log aggregation
- Audit trails

#### I. CI/CD & DevOps
- Pipeline configuration
- Build automation
- Test automation
- Deployment strategy
- Container security

#### J. Documentation
- Code comments and godoc
- README completeness
- API documentation
- Deployment guides
- Architecture documentation

#### K. Code Quality
- Go conventions and idioms
- Code duplication
- Function size and complexity
- Naming conventions
- Package dependencies

#### L. Performance
- Algorithm efficiency
- Resource usage
- Caching strategies
- Database query optimization
- Benchmark availability

### 3. Issue Creation

For each area of improvement identified:

**Issue Template**:
```markdown
## Description
[Clear description of the issue or improvement needed]

## Goals
- [Specific, measurable goal 1]
- [Specific, measurable goal 2]

## Files to modify
- `path/to/file1.go`
- `path/to/file2.go`

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Tests added/updated
- [ ] Documentation updated

---
**Estimated Effort**: [Small/Medium/Large]
**Priority**: [High/Medium/Low]
**Category**: [Testing/Security/Enhancement/etc.]
```

**Issue Guidelines**:
1. **Focused**: Each issue addresses one specific improvement
2. **Actionable**: Clear steps and acceptance criteria
3. **Scoped**: Can be completed by a single coding agent
4. **Testable**: Includes verification criteria
5. **Documented**: Explains why and how

### 4. Prioritization

**Priority Levels**:

- **High**: Security vulnerabilities, critical bugs, major quality issues
- **Medium**: Code quality improvements, feature enhancements, moderate technical debt
- **Low**: Nice-to-have improvements, optimizations, cosmetic changes

**Effort Estimation**:

- **Small**: 1-2 days (< 200 lines of code)
- **Medium**: 3-5 days (200-500 lines of code)
- **Large**: 1-2 weeks (> 500 lines of code)

### 5. Documentation

Create the following documents:

1. **README.md**: Navigation guide for the review
2. **CODE_REVIEW_SUMMARY.md**: Detailed findings and analysis
3. **CODE_REVIEW_ISSUES.md**: All issues with full details
4. **NEXT_STEPS.md**: Implementation roadmap
5. **QUICK_START_ISSUES.md**: Quick guide for creating GitHub issues

### 6. Roadmap Creation

Organize issues into implementation phases:

**Phase 1: Security & Critical** (Week 1)
- Close security gaps
- Fix critical bugs
- Address high-priority quality issues

**Phase 2: Testing & Quality** (Week 2)
- Increase test coverage
- Improve test quality
- Fix quality issues

**Phase 3: Code Quality** (Week 3)
- Refactor duplications
- Improve maintainability
- Update documentation

**Phase 4: Observability** (Week 4)
- Enhance monitoring
- Improve logging
- Add metrics

**Phase 5: Features & Enhancements** (Ongoing)
- New features
- Optimizations
- Nice-to-have improvements

## Best Practices

### Do's ✅

- **Be specific**: Provide exact file paths and line numbers where relevant
- **Be constructive**: Focus on improvements, not criticism
- **Be actionable**: Each issue should be implementable
- **Be thorough**: Cover all aspects of the codebase
- **Be objective**: Base findings on standards and best practices
- **Include examples**: Show code samples when helpful
- **Provide context**: Explain why something is an issue

### Don'ts ❌

- **Don't be vague**: Avoid general statements without specifics
- **Don't nitpick**: Focus on meaningful improvements
- **Don't assume**: Verify your findings by running tests
- **Don't ignore strengths**: Acknowledge what's working well
- **Don't create huge issues**: Keep issues focused and scoped
- **Don't skip verification**: Always validate your assessment

## Tools to Use

### Analysis Tools

```bash
# Test coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Code statistics
find . -name "*.go" ! -path "./vendor/*" | xargs wc -l

# Find TODOs and FIXMEs
grep -r "TODO\|FIXME" --include="*.go"

# Check for common issues
go vet ./...

# Security scanning (if available)
gosec ./...
```

### Validation

Always validate findings:
1. Run the full test suite
2. Check test coverage reports
3. Review actual code behavior
4. Verify security assumptions
5. Test critical paths manually

## Output Format

### Grading Scale

- **A (Excellent)**: Exceptional quality, minimal improvements needed
- **B+ (Very Good)**: Strong fundamentals, some enhancements recommended
- **B (Good)**: Solid codebase, several improvements needed
- **C (Acceptable)**: Functional but needs significant work
- **D (Poor)**: Major issues, substantial refactoring required

### Summary Structure

```markdown
## Overall Grade: [A/B+/B/C/D]

### Strengths
- [Strength 1]
- [Strength 2]

### Areas for Improvement
- [Area 1]
- [Area 2]

### Metrics
- Files analyzed: X
- Lines of code: Y
- Test coverage: Z%
- Issues identified: N
```

## Review Frequency

- **Major reviews**: Quarterly or before major releases
- **Focused reviews**: After significant features
- **Security reviews**: Before production deployment
- **Ad-hoc reviews**: When requested or after incidents

## Follow-up

After creating issues:
1. Create a GitHub Project or Milestone
2. Assign issues to team members or coding agents
3. Track implementation progress
4. Re-review after significant improvements
5. Update grades and metrics

---

## Example Review Checklist

Use this checklist to ensure comprehensive coverage:

- [ ] All source files analyzed
- [ ] Test coverage measured
- [ ] Security patterns reviewed
- [ ] Error handling examined
- [ ] API design evaluated
- [ ] Database patterns checked
- [ ] Configuration reviewed
- [ ] Observability assessed
- [ ] CI/CD pipelines examined
- [ ] Documentation verified
- [ ] Code quality evaluated
- [ ] Performance considerations noted
- [ ] Issues prioritized
- [ ] Roadmap created
- [ ] Documentation complete

---

**Version**: 1.0  
**Last Updated**: 2025-10-24  
**Maintained by**: Development Team
