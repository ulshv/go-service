package user

type UserService struct {
	userRepository *userRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepository: newUserRepository(),
	}
}

func (us *UserService) GetUserById(id int) (*User, error) {
	return us.userRepository.getUserById(id)
}

func (us *UserService) FindUserByEmail(username string) (*User, error) {
	return us.userRepository.findUserByEmail(username)
}

func (us *UserService) CreateUser(user User) (*User, error) {
	return us.userRepository.createUser(user)
}
