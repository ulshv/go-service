package user

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/ulshv/go-service/internal/logger"
)

type UserService struct {
	userRepository *userRepository
	logger         *slog.Logger
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		userRepository: newUserRepository(db),
		logger:         logger.NewLogger("UserService"),
	}
}

func (s *UserService) GetUserById(id int) (*User, error) {
	return s.userRepository.getUserById(id)
}

func (s *UserService) FindUserByEmail(username string) (*User, error) {
	return s.userRepository.findUserByEmail(username)
}

func (s *UserService) CreateUser(user User) (*User, error) {
	s.logger.Info("CreateUser", "email", user.Email)
	u, err := s.userRepository.createUser(user)
	if err != nil {
		// For PostgreSQL unique violation
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, ErrEmailTaken
			}
		}
		// For SQLite unique constraint
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, ErrEmailTaken
			}
		}
		return nil, err
	}
	return u, nil
}
