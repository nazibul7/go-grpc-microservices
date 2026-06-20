package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending migrations to the database.
func RunMigrations(dbURL string) error {

	// Create a migration instance.
	// Parameter 1: location of migration files.
	// Parameter 2: database connection URL.
	m, err := migrate.New("file://internal/migrations/", dbURL)
	if err != nil {
		return err
	}

	// Run all migrations that haven't been applied yet.
	// Ignore ErrNoChange because it means the database
	// is already at the latest migration version.
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("migrations applied")
	return nil
}
