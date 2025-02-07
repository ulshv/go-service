package user

import (
	"log/slog"

	"github.com/ulshv/online-store-app/backend-go/internal/logger"
)

type UserService struct {
	userRepository *userRepository
	logger         *slog.Logger
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: newUserRepository(),
		logger:         logger.NewLogger("UserService"),
	}
}

func (us *UserService) GetUserById(id int) (*User, error) {
	return us.userRepository.getUserById(id)
}

func (us *UserService) FindUserByEmail(username string) (*User, error) {
	return us.userRepository.findUserByEmail(username)
}

func (us *UserService) CreateUser(user User) (*User, error) {
	us.logger.Info("CreateUser", "email", user.Email)
	return us.userRepository.createUser(user)
}
