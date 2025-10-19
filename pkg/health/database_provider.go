package health

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// DatabaseHealthCheckProvider checks the health of the database connection.
type DatabaseHealthCheckProvider struct {
	db     *gorm.DB
	scopes []Scope
}

// NewDatabaseHealthCheckProvider creates a new database health check provider.
// The provider can be configured with multiple scopes.
func NewDatabaseHealthCheckProvider(db *gorm.DB, scopes ...Scope) *DatabaseHealthCheckProvider {
	// Default to startup and ready if no scopes provided
	if len(scopes) == 0 {
		scopes = []Scope{ScopeStartup, ScopeReady}
	}
	return &DatabaseHealthCheckProvider{
		db:     db,
		scopes: scopes,
	}
}

// Name returns the name of this health check.
func (d *DatabaseHealthCheckProvider) Name() string {
	return "database"
}

// Check executes the database health check.
func (d *DatabaseHealthCheckProvider) Check() (*CheckResult, error) {
	if d.db == nil {
		return &CheckResult{
			Status: StatusDown,
			Details: map[string]interface{}{
				"error": "database not initialized",
			},
		}, nil
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return &CheckResult{
			Status: StatusDown,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}, nil
	}

	// Ping database with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return &CheckResult{
			Status: StatusDown,
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}, nil
	}

	// Get database stats
	stats := sqlDB.Stats()
	return &CheckResult{
		Status: StatusUp,
		Details: map[string]interface{}{
			"max_open_connections": stats.MaxOpenConnections,
			"open_connections":     stats.OpenConnections,
			"in_use":               stats.InUse,
			"idle":                 stats.Idle,
		},
	}, nil
}

// Scopes returns the scopes for this health check.
func (d *DatabaseHealthCheckProvider) Scopes() []Scope {
	return d.scopes
}
