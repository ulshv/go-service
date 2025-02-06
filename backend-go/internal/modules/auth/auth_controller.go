package auth

import "net/http"

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
	w.Write([]byte("Login"))
}

func (c *AuthController) registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register"))
}

func (c *AuthController) logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout"))
}
