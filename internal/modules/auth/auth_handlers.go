package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ulshv/go-service/internal/core/httperrs"
	"github.com/ulshv/go-service/internal/modules/user"
	"github.com/ulshv/go-service/pkg/logs"
	"github.com/ulshv/go-service/pkg/utils/httputils"
)

type authHandlers struct {
	authService *authService
	logger      *slog.Logger
}

func newAuthHandlers(authService *authService) *authHandlers {
	return &authHandlers{
		authService: authService,
		logger:      logs.NewLogger("AuthHandlers"),
	}
}

func (h *authHandlers) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", h.registerHandler)
	mux.HandleFunc("POST /api/v1/auth/login", h.loginHandler)
	mux.HandleFunc("GET /api/v1/auth/me", h.meHandler)
}

func (h *authHandlers) registerHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("registerHandler")
	var registerDto RegisterDto
	err := httputils.DecodeBody(w, r, &registerDto)
	if err != nil {
		return
	}
	h.logger.Debug("registerHandler, parsed DTO", "email", registerDto.Email)
	result, err := h.authService.register(registerDto.Email, registerDto.Password)
	h.logger.Debug("after register", "result", result, "err", err)
	if err != nil {
		respStatus := http.StatusInternalServerError
		errCode := httperrs.ErrCodeUnknown
		if err == user.ErrEmailTaken {
			respStatus = http.StatusConflict
			errCode = httperrs.ErrEmailTaken
		}
		slog.Debug("received err, writing to client", "err", "err")
		httputils.WriteErrorJson(w, err, errCode, respStatus)
		return
	}
	h.logger.Debug("writing json response to client", "result", result)
	httputils.WriteJson(w, result)
}

func (h *authHandlers) loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginDto LoginDto
	err := httputils.DecodeBody(w, r, &loginDto)
	if err != nil {
		return
	}
	result, err := h.authService.login(loginDto.Email, loginDto.Password)
	if err != nil {
		errCode := httperrs.ErrCodeUnknown
		if errors.Is(err, errInvalidEmailOrPassword) {
			errCode = httperrs.ErrInvalidEmailOrPassword
		}
		httputils.WriteErrorJson(w, err, errCode, http.StatusUnauthorized)
		return
	}
	httputils.WriteJson(w, result)
}

func (h *authHandlers) meHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented yet."))
}
