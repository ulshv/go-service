package server

import (
	"fmt"
	"net/http"

	"github.com/ulshv/go-service/internal/core/application"
)

func NewServer(address, port string, app *application.App) *http.Server {
	mux := http.NewServeMux()

	registerHandlers(mux, app)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", address, port),
		Handler: mux,
	}
	return server
}
