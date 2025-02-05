package main

import (
	"log/slog"

	"github.com/ulshv/online-store-app/backend-go/server"
)

func main() {
	slog.Info("Starting the server...")

	s := server.NewServer(server.NewServerOptions{
		Address: "localhost",
		Port:    "5000",
	})

	s.Listen()
}
