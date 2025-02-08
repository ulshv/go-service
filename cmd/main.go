package main

import (
	"log"
	"log/slog"

	"github.com/ulshv/go-service/internal/application"
	"github.com/ulshv/go-service/internal/database"
	"github.com/ulshv/go-service/internal/server"
	"github.com/ulshv/go-service/internal/utils/envutils"
)

func main() {
	envutils.LoadEnvFiles(".env")
	dbConfig := database.Config{
		Host:     envutils.RequireEnv("POSTGRES_HOST"),
		Port:     envutils.RequireEnv("POSTGRES_PORT"),
		User:     envutils.RequireEnv("POSTGRES_USER"),
		Password: envutils.RequireEnv("POSTGRES_PASSWORD"),
		DBName:   envutils.RequireEnv("POSTGRES_DB"),
	}
	slog.Info("Initializing the application")
	app, err := application.NewApp(dbConfig)
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}
	defer app.Close()
	srv := server.NewServer("0.0.0.0", "5000", app)
	slog.Info("Starting the listener", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
