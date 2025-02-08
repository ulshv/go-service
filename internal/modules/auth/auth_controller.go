package auth

import (
	"log/slog"
	"net/http"

	"github.com/ulshv/online-store-app/backend-go/internal/logger"
	"github.com/ulshv/online-store-app/backend-go/internal/modules/user"
	"github.com/ulshv/online-store-app/backend-go/internal/utils/httputils"
)

type authController struct {
	authService *authService
	logger      *slog.Logger
}

func newAuthController(authService *authService) *authController {
	return &authController{
		authService: authService,
		logger:      logger.NewLogger("AuthController"),
	}
}

func (c *authController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", c.registerHandler)
	mux.HandleFunc("POST /api/v1/auth/login", c.loginHandler)
}

func (c *authController) registerHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("registerHandler")
	var registerDto registerDto
	err := httputils.DecodeBody(w, r, &registerDto)
	if err != nil {
		return
	}
	c.logger.Debug("registerHandler, parsed DTO", "email", registerDto.Email)
	result, err := c.authService.register(registerDto.Email, registerDto.Password)
	c.logger.Debug("after register", "result", result, "err", err)
	if err != nil {
		respStatus := http.StatusInternalServerError
		if err == user.ErrEmailTaken {
			respStatus = http.StatusConflict
		}
		slog.Debug("received err, writing to client", "err", "err")
		httputils.WriteErrorJson(w, err.Error(), respStatus)
		return
	}
	c.logger.Debug("writing json response to client", "result", result)
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
