package mw

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ulshv/go-service/internal/core/httperrs"
	"github.com/ulshv/go-service/pkg/logs"
	"github.com/ulshv/go-service/pkg/utils/httputils"
	"github.com/ulshv/go-service/pkg/utils/jwtutils"
)

var jwt = jwtutils.NewJWT()
var logger = logs.NewLogger("middlewares")

type mwKeyTypes int

const (
	userIdKey mwKeyTypes = iota + 1
	authErrKey
)

func Chain(handlers ...http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r
		for _, h := range handlers {
			h(w, req)
			hasErr, ok := r.Context().Value(authErrKey).(bool)
			if ok && hasErr {
				break
			}
		}
	}
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		return
	}
	jwtToken := strings.TrimPrefix(accessToken, "Bearer ")
	claims, err := jwt.ValidateAccessToken(jwtToken)
	if err != nil {
		logger.Debug("ValidateAccessToken error", "error", err)
		if errors.Is(err, jwtutils.ErrInvalidToken) {
			httputils.WriteErrorJson(w, err, httperrs.ErrCodeAccessTokenInvalid, http.StatusUnauthorized)
			setAuthErr(r)
		}
		if errors.Is(err, jwtutils.ErrTokenExpired) {
			httputils.WriteErrorJson(w, err, httperrs.ErrCodeAccessTokenExpired, http.StatusUnauthorized)
			setAuthErr(r)
		}
		// allow passing the MW if JWT cannot be parsed, if auth is needed - it will be catched with `AuthRequired` MW
		return
	}
	userId := claims.UserId
	if userId != 0 {
		ctx := context.WithValue(r.Context(), userIdKey, userId)
		*r = *r.WithContext(ctx)
	}
}

func AuthRequired(w http.ResponseWriter, r *http.Request) {
	_, ok := GetUserId(r)
	if !ok {
		setAuthErr(r)
		httputils.WriteErrorJson(w, httperrs.ErrUnauthorized, httperrs.ErrCodeUnautorized, http.StatusUnauthorized)
	}
}

func GetUserId(r *http.Request) (int, bool) {
	userId, ok := r.Context().Value(userIdKey).(int)
	return userId, ok
}

func setAuthErr(r *http.Request) {
	ctx := context.WithValue(r.Context(), authErrKey, true)
	*r = *r.WithContext(ctx)
}
