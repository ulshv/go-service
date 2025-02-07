package main

import (
	"log"
	"log/slog"

	"github.com/ulshv/online-store-app/backend-go/internal/application"
	"github.com/ulshv/online-store-app/backend-go/internal/server"
)

func main() {
	slog.Info("Initializing the application")

	app := application.NewApp()
	srv := server.NewServer("0.0.0.0", "5000", app)

	slog.Info("Starting the listener", "address", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
