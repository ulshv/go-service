package migrations

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	db_mod "github.com/ulshv/go-service/internal/database"
)

func RunMigrations(db *sqlx.DB, migrationsPath string, logger *slog.Logger, dbType db_mod.DBType) error {
	logger.Info("RunMigrations", "database_type", dbType)

	// Get current working directory and create absolute path to migrations
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("could not get working directory", "error", err)
		return fmt.Errorf("could not get working directory: %w", err)
	}
	absolutePath := filepath.Join(cwd, "../../../", migrationsPath)
	sourceURL := fmt.Sprintf("file://%s", absolutePath)
	logger.Debug("migrations source URL", "url", sourceURL)

	var driver database.Driver
	switch dbType {
	case db_mod.PostgreSQL:
		driver, err = postgres.WithInstance(db.DB, &postgres.Config{})
	case db_mod.SQLite:
		driver, err = sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	default:
		logger.Error("unsupported database type for migrations", "database_type", dbType)
		return fmt.Errorf("unsupported database type for migrations: %s", dbType)
	}
	if err != nil {
		logger.Error("could not create migrations driver", "error", err)
		return fmt.Errorf("could not create migrations driver: %w", err)
	}
	logger.Debug("creating migrate instance")
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		string(dbType),
		driver,
	)
	if err != nil {
		logger.Error("could not create migrate instance", "error", err)
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	logger.Info("Running migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}
	logger.Info("Database migrations completed successfully")
	return nil
}
