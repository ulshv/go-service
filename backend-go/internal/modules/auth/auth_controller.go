package auth

import (
	"log/slog"
	"net/http"

	"github.com/ulshv/online-store-app/backend-go/internal/utils/httputils"
)

type authController struct {
	authService *authService
}

func newAuthController(authService *authService) *authController {
	return &authController{authService: authService}
}

func (c *authController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", c.registerHandler)
	mux.HandleFunc("POST /api/v1/auth/login", c.loginHandler)
}

func (c *authController) registerHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("registerHandler")
	var registerDto registerDto
	err := httputils.DecodeBody(w, r, &registerDto)
	if err != nil {
		return
	}
	slog.Info("registerHandler, parsed DTO", "registerDto", registerDto)
	result, err := c.authService.register(registerDto.Email, registerDto.Password)
	if err != nil {
		slog.Info("registerHandler, error", "error", err)
		httputils.WriteErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("registerHandler, result", "result", result)
	httputils.WriteJson(w, result)
}

func (c *authController) loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginDto loginDto
	err := httputils.DecodeBody(w, r, &loginDto)
	if err != nil {
		return
	}
	token, err := c.authService.login(loginDto.Email, loginDto.Password)
	if err != nil {
		httputils.WriteErrorJson(w, err.Error(), http.StatusUnauthorized)
		return
	}
	httputils.WriteJson(w, map[string]string{"token": token})
}

func (c *authController) logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout"))
}
