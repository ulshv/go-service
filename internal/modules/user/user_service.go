package user

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/ulshv/go-service/pkg/logs"
)

type UserService struct {
	userRepository *userRepository
	logger         *slog.Logger
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		userRepository: newUserRepository(db),
		logger:         logs.NewLogger("UserService"),
	}
}

func (s *UserService) GetUserByID(id int) (*User, error) {
	return s.userRepository.getUserByID(id)
}

func (s *UserService) FindUserByEmail(username string) (*User, error) {
	return s.userRepository.findUserByEmail(username)
}

func (s UserService) NewUser(email, passwordHash string) User {
	return User{
		Email:        email,
		PasswordHash: passwordHash,
	}
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
