package server

import (
	"net/http"

	"github.com/ulshv/online-store-app/backend-go/internal/application"
)

func registerRoutes(mux *http.ServeMux, app *application.App) *http.ServeMux {
	app.AuthModule.AuthController.RegisterRoutes(mux)

	return mux
}
