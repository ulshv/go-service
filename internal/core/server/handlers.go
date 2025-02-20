package server

import (
	"net/http"

	"github.com/ulshv/go-service/internal/core/application"
)

func registerHandlers(mux *http.ServeMux, app *application.App) *http.ServeMux {
	app.AuthModule.AuthHandlers.RegisterHandlers(mux)

	return mux
}
