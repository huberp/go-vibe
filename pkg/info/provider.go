package info

// InfoProvider defines the interface for providing information to the `/info` endpoint.
// Implementations of this interface can provide any type of application information
// such as build details, statistics, or custom metrics.
type InfoProvider interface {
	// Name returns the unique name of this provider.
	// This name will be used as the key in the aggregated info response.
	Name() string

	// Info returns the information provided by this provider.
	// The map can contain any JSON-serializable data.
	// Returns an error if the information cannot be retrieved.
	Info() (map[string]any, error)
}
