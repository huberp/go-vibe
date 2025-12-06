package info

import "sync"

// Registry manages a collection of InfoProvider instances.
// It provides thread-safe registration and aggregation of information providers.
type Registry struct {
	mu        sync.RWMutex
	providers []InfoProvider
}

// NewRegistry creates a new InfoProvider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make([]InfoProvider, 0),
	}
}

// Register adds a new InfoProvider to the registry.
// This method is thread-safe and can be called from multiple goroutines.
func (r *Registry) Register(provider InfoProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = append(r.providers, provider)
}

// GetAll aggregates information from all registered providers.
// Returns a map where keys are provider names and values are the information maps.
// If a provider returns an error, that provider's data is omitted from the result.
func (r *Registry) GetAll() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]any)
	for _, provider := range r.providers {
		info, err := provider.Info()
		if err == nil {
			result[provider.Name()] = info
		}
	}
	return result
}
