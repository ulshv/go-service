package main

import (
	"log/slog"

	"github.com/ulshv/go-web-app/internal/server"
)

func main() {
	slog.Info("Starting go-web-app...")

	s := server.NewServer(server.NewServerOptions{
		Address: "localhost",
		Port:    "5000",
	})

	s.Listen()
}
