package health

import "sync"

// Registry manages a collection of HealthCheckProvider instances.
// It provides thread-safe registration and aggregation of health check providers.
type Registry struct {
	mu        sync.RWMutex
	providers []HealthCheckProvider
}

// NewRegistry creates a new HealthCheckProvider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make([]HealthCheckProvider, 0),
	}
}

// Register adds a new HealthCheckProvider to the registry.
// This method is thread-safe and can be called from multiple goroutines.
func (r *Registry) Register(provider HealthCheckProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = append(r.providers, provider)
}

// Check executes all health checks that match the given scope.
// Returns a map where keys are provider names and values are the check results.
// If a provider returns an error, that provider's data is omitted from the result.
// When scope is empty (nil), it aggregates all providers but ensures each provider
// is only checked once even if it appears in multiple scopes.
func (r *Registry) Check(scope *Scope) map[string]*CheckResult {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*CheckResult)
	checked := make(map[string]bool) // Track which providers we've already checked

	for _, provider := range r.providers {
		// Skip if we've already checked this provider (for /health endpoint)
		if scope == nil {
			if checked[provider.Name()] {
				continue
			}
		} else {
			// Check if provider has the requested scope
			hasScope := false
			for _, s := range provider.Scopes() {
				if s == *scope {
					hasScope = true
					break
				}
			}
			if !hasScope {
				continue
			}
		}

		checkResult, err := provider.Check()
		if err == nil && checkResult != nil {
			result[provider.Name()] = checkResult
			checked[provider.Name()] = true
		}
	}

	return result
}

// OverallStatus determines the overall health status based on all component statuses.
// Returns StatusDown if any component is down, otherwise StatusUp.
func OverallStatus(components map[string]*CheckResult) Status {
	if len(components) == 0 {
		return StatusUp
	}

	for _, component := range components {
		if component.Status == StatusDown {
			return StatusDown
		}
	}

	return StatusUp
}
