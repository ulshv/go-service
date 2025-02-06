package user

import "sync"

type userRepository struct {
	users map[int]*User

	lock sync.RWMutex
}

func newUserRepository() *userRepository {
	return &userRepository{
		users: make(map[int]*User),
		lock:  sync.RWMutex{},
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
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, user := range r.users {
		if user.Email == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *userRepository) createUser(user User) (*User, error) {
	existingUser, err := r.findUserByEmail(user.Email)
	if err == nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	user.Id = len(r.users) + 1
	r.users[user.Id] = &user
	return &user, nil
}
