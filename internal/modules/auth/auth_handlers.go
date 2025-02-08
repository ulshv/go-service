package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/ulshv/go-service/internal/logger"
	"github.com/ulshv/go-service/internal/modules/user"
	"github.com/ulshv/go-service/internal/utils/httputils"
)

type authHandlers struct {
	authService *authService
	logger      *slog.Logger
}

func newAuthHandlers(authService *authService) *authHandlers {
	return &authHandlers{
		authService: authService,
		logger:      logger.NewLogger("AuthHandlers"),
	}
}

func (h *authHandlers) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", h.registerHandler)
	mux.HandleFunc("POST /api/v1/auth/login", h.loginHandler)
	mux.HandleFunc("GET /api/v1/auth/me", h.meHandler)
}

func (h *authHandlers) registerHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("registerHandler")
	var registerDto registerDto
	err := httputils.DecodeBody(w, r, &registerDto)
	if err != nil {
		return
	}
	h.logger.Debug("registerHandler, parsed DTO", "email", registerDto.Email)
	result, err := h.authService.register(registerDto.Email, registerDto.Password)
	h.logger.Debug("after register", "result", result, "err", err)
	if err != nil {
		respStatus := http.StatusInternalServerError
		if err == user.ErrEmailTaken {
			respStatus = http.StatusConflict
		}
		slog.Debug("received err, writing to client", "err", "err")
		httputils.WriteErrorJson(w, err.Error(), respStatus)
		return
	}
	h.logger.Debug("writing json response to client", "result", result)
	httputils.WriteJson(w, result)
}

func (h *authHandlers) loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginDto loginDto
	err := httputils.DecodeBody(w, r, &loginDto)
	if err != nil {
		return
	}
	result, err := h.authService.login(loginDto.Email, loginDto.Password)
	if err != nil {
		httputils.WriteErrorJson(w, err.Error(), http.StatusUnauthorized)
		return
	}
	httputils.WriteJson(w, result)
}

func (h *authHandlers) meHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("meHandler")
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		httputils.WriteErrorJson(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	h.logger.Debug("meHandler - got access token", "token", accessToken)
	claims, err := h.authService.jwt.ValidateAccessToken(accessToken)
	if err != nil {
		httputils.WriteErrorJson(w, err.Error(), http.StatusUnauthorized)
		return
	}
	h.logger.Debug("meHandler - parsed token", "claims", claims)
	user, err := h.authService.userService.GetUserById(claims.UserId)
	if err != nil {
		httputils.WriteErrorJson(w, "invalid token", http.StatusUnauthorized)
		return
	}
	user.PasswordHash = ""
	httputils.WriteJson(w, user)
}
