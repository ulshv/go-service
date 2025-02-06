package auth

import (
	"errors"
)

var (
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

type loginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
