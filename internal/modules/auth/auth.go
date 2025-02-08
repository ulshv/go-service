package auth

import (
	"errors"

	"github.com/ulshv/go-service/internal/utils/jwtutils"
)

var (
	errInvalidEmailOrPassword = errors.New("invalid email or password")
)

type registerDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResultDto struct {
	UserId int                `json:"user_id"`
	Tokens jwtutils.TokenPair `json:"tokens"`
}

type loginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResultDto struct {
	Tokens jwtutils.TokenPair `json:"tokens"`
}
