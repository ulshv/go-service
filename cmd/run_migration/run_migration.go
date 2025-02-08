package main

import (
	"os"

	"github.com/ulshv/online-store-app/backend-go/internal/database"
	"github.com/ulshv/online-store-app/backend-go/internal/database/migrations"
	"github.com/ulshv/online-store-app/backend-go/internal/logger"
	"github.com/ulshv/online-store-app/backend-go/internal/utils/envutils"
)

func main() {
	logger := logger.NewLogger("run_migration")
	logger.Info("Starting the migration")

	envutils.LoadEnvFiles(".env")

	db, err := database.NewConnection(database.Config{
		Type:     database.PostgreSQL,
		Host:     envutils.RequireEnv("POSTGRES_HOST"),
		Port:     envutils.RequireEnv("POSTGRES_PORT"),
		User:     envutils.RequireEnv("POSTGRES_USER"),
		Password: envutils.RequireEnv("POSTGRES_PASSWORD"),
		DBName:   envutils.RequireEnv("POSTGRES_DB"),
	})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	if err = migrations.RunMigrations(db, "migrations", logger, database.PostgreSQL); err != nil {
		logger.Error("Failed to run database migrations", "error", err)
		os.Exit(1)
	}

	logger.Info("Database migrations completed successfully")
}
