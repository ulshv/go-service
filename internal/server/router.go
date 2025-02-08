package server

import (
	"net/http"

	"github.com/ulshv/go-service/internal/application"
)

func registerRoutes(mux *http.ServeMux, app *application.App) *http.ServeMux {
	app.AuthModule.AuthController.RegisterRoutes(mux)

	return mux
}
