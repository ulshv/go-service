package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ulshv/go-service/internal/logger"
)

type DBType string

const (
	PostgreSQL DBType = "postgres"
	SQLite     DBType = "sqlite3"
)

type Config struct {
	Type     DBType
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewConnection(cfg Config) (*sqlx.DB, error) {
	logger := logger.NewLogger("Database")

	logger.Info("Connecting to database...", "config", cfg)

	var dsn string
	switch cfg.Type {
	case PostgreSQL:
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	case SQLite:
		dsn = cfg.DBName
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	db, err := sqlx.Connect(string(cfg.Type), dsn)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Successfully connected to database")

	if err := db.Ping(); err != nil {
		logger.Error("failed to ping database", "error", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to database", "type", cfg.Type)
	return db, nil
}
