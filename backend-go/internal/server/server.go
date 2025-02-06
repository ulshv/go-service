package server

import (
	"fmt"
	"net/http"

	"github.com/ulshv/online-store-app/backend-go/internal/application"
)

func NewServer(address, port string, app *application.App) *http.Server {
	mux := http.NewServeMux()

	registerRoutes(mux, app)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", address, port),
		Handler: mux,
	}
	return server
}
