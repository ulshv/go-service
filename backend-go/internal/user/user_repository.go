package user

import "sync"

type UserRepository struct {
	users map[int]*User

	lock sync.RWMutex
}

func (r *UserRepository) GetUserById(id int) (*User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *UserRepository) FindUserByName(username string) (*User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *UserRepository) CreateUser(user User) (*User, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	u, err := r.FindUserByName(user.Username)
	if err == nil {
		return nil, err
	}
	if u != nil {
		return nil, ErrUserExists
	}
	u.Id = len(r.users) + 1
	r.users[u.Id] = &user
	return u, nil
}
