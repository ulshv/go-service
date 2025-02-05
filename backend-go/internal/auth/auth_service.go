package auth

import "github.com/ulshv/online-store-app/backend-go/internal/user"

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(
