package auth

import (
	"errors"

	"github.com/ulshv/go-service/pkg/utils/jwtutils"
)

var (
	errInvalidEmailOrPassword = errors.New("invalid email or password")
)

type RegisterDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResultDto struct {
	UserId int                `json:"user_id"`
	Tokens jwtutils.TokenPair `json:"tokens"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResultDto struct {
	Tokens jwtutils.TokenPair `json:"tokens"`
}
