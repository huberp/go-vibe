package health

import "time"

// SimpleHealthCheckProvider is a basic example health check provider.
// This demonstrates how to create custom health checks for your application.
// It always returns UP status and can be used as a template for real implementations.
type SimpleHealthCheckProvider struct {
	name   string
	scopes []Scope
}

// NewSimpleHealthCheckProvider creates a new simple health check provider.
// This is useful for testing or as a basic liveness check.
func NewSimpleHealthCheckProvider(name string, scopes ...Scope) *SimpleHealthCheckProvider {
	// Default to liveness scope if none provided
	if len(scopes) == 0 {
		scopes = []Scope{ScopeLive}
	}
	return &SimpleHealthCheckProvider{
		name:   name,
		scopes: scopes,
	}
}

// Name returns the name of this health check.
func (s *SimpleHealthCheckProvider) Name() string {
	return s.name
}

// Check executes the health check.
// This implementation always returns UP status with the current timestamp.
func (s *SimpleHealthCheckProvider) Check() (*CheckResult, error) {
	return &CheckResult{
		Status: StatusUp,
		Details: map[string]interface{}{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	}, nil
}

// Scopes returns the scopes for this health check.
func (s *SimpleHealthCheckProvider) Scopes() []Scope {
	return s.scopes
}
