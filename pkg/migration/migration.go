package migration

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

// RunMigrations runs database migrations
func RunMigrations(databaseURL string, logger *zap.Logger) error {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Migrations applied successfully")
	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(databaseURL string, logger *zap.Logger) error {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Rollback one step
	if err := m.Steps(-1); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	logger.Info("Migration rolled back successfully")
	return nil
}

// MigrateToVersion migrates to a specific version
func MigrateToVersion(databaseURL string, version uint, logger *zap.Logger) error {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Migrate(version); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("Already at target version", zap.Uint("version", version))
			return nil
		}
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	logger.Info("Migrated to version successfully", zap.Uint("version", version))
	return nil
}
