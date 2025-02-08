package auth

import (
	"log/slog"

	"github.com/ulshv/go-service/internal/logger"
	"github.com/ulshv/go-service/internal/modules/user"
	"github.com/ulshv/go-service/internal/utils/jwtutils"
)

type authService struct {
	userService *user.UserService
	logger      *slog.Logger
	jwt         *jwtutils.Jwt
}

func newAuthService(
	userService *user.UserService,
) *authService {
	return &authService{
		userService: userService,
		logger:      logger.NewLogger("AuthService"),
		jwt:         jwtutils.NewJWT(),
	}
}

func (s *authService) register(email, password string) (*registerResultDto, error) {
	s.logger.Info("register", "email", email)
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	payload := s.userService.NewUser(email, passwordHash)
	s.logger.Debug("register - created payload, now trying to create user", "payload", payload)
	user, err := s.userService.CreateUser(payload)
	s.logger.Debug("register - CreateUser result", "user", user, "err", err)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("register - now generate toekn")
	tokenPair, err := s.jwt.GenerateTokenPair(user.Id)
	if err != nil {
		return nil, err
	}
	s.logger.Debug("register - generated token")
	return &registerResultDto{
		UserId: user.Id,
		Tokens: tokenPair,
	}, nil
}

func (s *authService) login(email, password string) (*loginResultDto, error) {
	s.logger.Info("login", "email", email)
	user, err := s.userService.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if !validatePassword(password, user.PasswordHash) {
		s.logger.Debug("login - invalid email or password")
		return nil, errInvalidEmailOrPassword
	}
	tokenPair, err := s.jwt.GenerateTokenPair(user.Id)
	if err != nil {
		return nil, err
	}
	return &loginResultDto{
		Tokens: tokenPair,
	}, nil
}
