package user

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/ulshv/go-service/pkg/logs"
)

type userRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func newUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db:     db,
		logger: logs.NewLogger("UserRepository"),
	}
}

func (r *userRepository) getUserById(id int) (*User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		r.logger.Error("failed to get user by id", "error", err)
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *userRepository) findUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		r.logger.Debug("user not found by email", "email", email)
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *userRepository) createUser(user User) (*User, error) {
	r.logger.Info("creating user", "email", user.Email)

	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, email, password_hash, created_at, updated_at`

	err := r.db.QueryRowx(
		query,
		user.Email,
		user.PasswordHash,
	).StructScan(&user)

	if err != nil {
		r.logger.Error("failed to create user", "error", err)
		return nil, err
	}

	return &user, nil
}
