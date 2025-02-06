package auth

import (
	"log/slog"

	"github.com/ulshv/online-store-app/backend-go/internal/modules/user"
)

type authService struct {
	userService *user.UserService
}

func newAuthService(
	userService *user.UserService,
) *authService {
	return &authService{
		userService: userService,
	}
}

func (s *authService) register(email, password string) (*registerResultDto, error) {
	slog.Info("register")
	payload := user.User{
		Email:        email,
		PasswordHash: hashPassword(password),
	}
	slog.Info("CreateUser")
	user, err := s.userService.CreateUser(payload)
	if err != nil {
		return nil, err
	}
	slog.Info("CreatedUser successfully")
	return &registerResultDto{
		UserId: user.Id,
		Token:  generateToken(user.Id),
	}, nil
}

func (s *authService) login(email, password string) (string, error) {
	user, err := s.userService.FindUserByEmail(email)
	if err != nil {
		return "", err
	}
	if validatePassword(password, user.PasswordHash) {
		return "", errInvalidEmailOrPassword
	}
	return generateToken(user.Id), nil
}
