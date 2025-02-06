package auth

import (
	"net/http"

	"github.com/ulshv/online-store-app/backend-go/internal/httputils"
)

type AuthController struct {
	authService *AuthService
}

func newAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/auth/login", c.loginHandler)
	mux.HandleFunc("/api/v1/auth/register", c.registerHandler)
	mux.HandleFunc("/api/v1/auth/logout", c.logoutHandler)
}

func (c *AuthController) loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginDto loginDto
	err := httputils.DecodeJsonAndHandleErr(w, r, &loginDto)
	if err != nil {
		return
	}
	token, err := c.authService.Login(loginDto.Email, loginDto.Password)
	if err != nil {
		httputils.ErrorJson(w, err.Error(), http.StatusUnauthorized)
		return
	}
	httputils.WriteJson(w, map[string]string{"token": token})
}

func (c *AuthController) registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register"))
}

func (c *AuthController) logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout"))
}
