package user

import (
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"username"`
	Name     string `json:"name"`
	Password string
}
