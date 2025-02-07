package user

import (
	"log/slog"
	"sync"

	"github.com/ulshv/online-store-app/backend-go/internal/logger"
)

type userRepository struct {
	users map[int]*User

	logger *slog.Logger
	lock   sync.RWMutex
}

func newUserRepository() *userRepository {
	return &userRepository{
		users:  make(map[int]*User),
		logger: logger.NewLogger("UserRepository"),
		lock:   sync.RWMutex{},
	}
}

func (r *userRepository) getUserById(id int) (*User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *userRepository) findUserByEmail(username string) (*User, error) {
	r.logger.Debug("findUserByEmail - locking for read")
	r.lock.RLock()
	defer r.lock.RUnlock()
	r.logger.Debug("findUserByEmail - locked for read")
	for _, user := range r.users {
		if user.Email == username {
			r.logger.Debug("findUserByEmail - found user", "user", user)
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *userRepository) createUser(user User) (*User, error) {
	r.logger.Info("createUser", "user", user)
	r.logger.Debug("createUser - check if user exists", "email", user.Email)
	existingUser, err := r.findUserByEmail(user.Email)
	r.logger.Debug("createUser - checked user", "existingUser", existingUser, "err", err)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		r.logger.Debug("createUser - user exists, returning ErrUserExists")
		return nil, ErrUserExists
	}
	r.logger.Debug("createUser - locking before creating user in the map")
	r.lock.Lock()
	defer r.lock.Unlock()
	user.Id = len(r.users) + 1
	r.users[user.Id] = &user
	r.logger.Debug("createUser - user added to the map", "user", user)
	return &user, nil
}
