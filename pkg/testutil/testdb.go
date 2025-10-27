package testutil

import (
	"database/sql"
	"myapp/internal/models"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDB provides database utilities for testing
type TestDB struct {
	DB       *gorm.DB
	fixtures *testfixtures.Loader
	sqlDB    *sql.DB
}

// NewTestDB creates a new test database with SQLite in-memory
func NewTestDB(t *testing.T) *TestDB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Auto-migrate the User model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get sql.DB from gorm: %v", err)
	}

	return &TestDB{
		DB:    db,
		sqlDB: sqlDB,
	}
}

// NewTestDBWithPostgres creates a new test database with PostgreSQL
// This requires a running PostgreSQL instance (use for integration tests)
func NewTestDBWithPostgres(t *testing.T, dsn string) *TestDB {
	t.Helper()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to setup PostgreSQL test database: %v", err)
	}

	// Auto-migrate the User model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get sql.DB from gorm: %v", err)
	}

	return &TestDB{
		DB:    db,
		sqlDB: sqlDB,
	}
}

// LoadFixtures loads test fixtures from the specified directory
func (tdb *TestDB) LoadFixtures(t *testing.T, fixturesDir string) {
	t.Helper()

	fixtures, err := testfixtures.New(
		testfixtures.Database(tdb.sqlDB),
		testfixtures.Dialect("sqlite3"),
		testfixtures.Directory(fixturesDir),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		t.Fatalf("Failed to create fixtures loader: %v", err)
	}

	tdb.fixtures = fixtures

	if err := fixtures.Load(); err != nil {
		t.Fatalf("Failed to load fixtures: %v", err)
	}
}

// LoadFixturesPostgres loads test fixtures for PostgreSQL
func (tdb *TestDB) LoadFixturesPostgres(t *testing.T, fixturesDir string) {
	t.Helper()

	fixtures, err := testfixtures.New(
		testfixtures.Database(tdb.sqlDB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(fixturesDir),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		t.Fatalf("Failed to create fixtures loader: %v", err)
	}

	tdb.fixtures = fixtures

	if err := fixtures.Load(); err != nil {
		t.Fatalf("Failed to load fixtures: %v", err)
	}
}

// ReloadFixtures reloads the fixtures (useful for tests that modify data)
func (tdb *TestDB) ReloadFixtures(t *testing.T) {
	t.Helper()

	if tdb.fixtures == nil {
		t.Fatal("Fixtures not initialized. Call LoadFixtures first.")
	}

	if err := tdb.fixtures.Load(); err != nil {
		t.Fatalf("Failed to reload fixtures: %v", err)
	}
}

// Cleanup cleans up the test database
func (tdb *TestDB) Cleanup() {
	if tdb.sqlDB != nil {
		tdb.sqlDB.Close()
	}
}

// TruncateTables truncates all tables in the database
func (tdb *TestDB) TruncateTables(t *testing.T) {
	t.Helper()

	// Delete all users
	if err := tdb.DB.Exec("DELETE FROM users").Error; err != nil {
		t.Fatalf("Failed to truncate users table: %v", err)
	}
}

// BeginTransaction starts a database transaction for test isolation
func (tdb *TestDB) BeginTransaction() *gorm.DB {
	return tdb.DB.Begin()
}
