package auth

import (
	"errors"
)

var (
	errInvalidEmailOrPassword = errors.New("invalid email or password")
)

type registerDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResultDto struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type loginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
