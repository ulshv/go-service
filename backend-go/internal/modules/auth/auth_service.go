package auth

import (
	"fmt"

	"github.com/ulshv/online-store-app/backend-go/internal/modules/user"
)

type AuthService struct {
	userService *user.UserService
}

func newAuthService(
	userService *user.UserService,
) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userService.FindUserByEmail(email)
	if err != nil {
		return "", err
	}
	// TODO hash password
	if user.Password != password {
		return "", ErrInvalidEmailOrPassword
	}
	// TODO generate token
	return fmt.Sprintf("token-%v", user.Id), nil
}

func (s *AuthService) Register(email, password string) (string, error) {
	// TODO hash password
	user := user.User{
		Email:    email,
		Password: password,
	}
	_, err := s.userService.CreateUser(user)
	if err != nil {
		return "", err
	}
	// TODO generate token
	return fmt.Sprintf("token-%v", user.Id), nil
}

func (s *AuthService) Logout(token string) error {
	// TODO invalidate token
	return nil
}
