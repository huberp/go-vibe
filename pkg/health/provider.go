package health

// Scope represents the scope in which a health check should be executed
type Scope string

const (
	// ScopeBase indicates the health check should only appear in /health
	ScopeBase Scope = "base"
	// ScopeStartup indicates the health check should appear in /health/startup, /health/readiness, and /health
	ScopeStartup Scope = "startup"
	// ScopeReady indicates the health check should appear in /health/readiness and /health
	ScopeReady Scope = "ready"
	// ScopeLive indicates the health check should appear in /health/liveness and /health
	ScopeLive Scope = "live"
)

// Status represents the health status of a component
type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Status  Status         `json:"status"`
	Details map[string]any `json:"details,omitempty"`
}

// ComponentHealth represents the health of a single component
type ComponentHealth struct {
	Status  Status         `json:"status"`
	Details map[string]any `json:"details,omitempty"`
}

// Response represents the overall health response
type Response struct {
	Status     Status                     `json:"status"`
	Components map[string]ComponentHealth `json:"components,omitempty"`
}

// HealthCheckProvider defines the interface for providing health check information.
// Implementations can check various aspects of application health such as
// database connectivity, external service availability, or custom metrics.
type HealthCheckProvider interface {
	// Name returns the unique name of this health check.
	// This name will be used as the key in the health response.
	Name() string

	// Check executes the health check and returns the result.
	// Returns an error if the health check cannot be performed.
	Check() (*CheckResult, error)

	// Scopes returns the scopes in which this health check should be executed.
	// A health check can be registered in multiple scopes.
	Scopes() []Scope
}
