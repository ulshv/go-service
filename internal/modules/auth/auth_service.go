package auth

import (
	"log/slog"

	"github.com/ulshv/go-service/internal/logger"
	"github.com/ulshv/go-service/internal/modules/user"
)

type authService struct {
	userService *user.UserService
	logger      *slog.Logger
}

func newAuthService(
	userService *user.UserService,
) *authService {
	return &authService{
		userService: userService,
		logger:      logger.NewLogger("AuthService"),
	}
}

func (s *authService) register(email, password string) (*registerResultDto, error) {
	s.logger.Info("register", "email", email)
	payload := user.User{
		Email:        email,
		PasswordHash: hashPassword(password),
	}
	s.logger.Debug("register - created payload, now trying to create user", "payload", payload)
	user, err := s.userService.CreateUser(payload)
	s.logger.Debug("register - CreateUser result", "user", user, "err", err)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("register - now generate toekn")
	token := generateToken(user.Id)
	s.logger.Debug("register - generated token")
	return &registerResultDto{
		UserId: user.Id,
		Token:  token,
	}, nil
}

func (s *authService) login(email, password string) (string, error) {
	s.logger.Info("login", "email", email)
	user, err := s.userService.FindUserByEmail(email)
	if err != nil {
		return "", err
	}
	if !validatePassword(password, user.PasswordHash) {
		s.logger.Debug("login - invalid email or password")
		return "", errInvalidEmailOrPassword
	}
	return generateToken(user.Id), nil
}
