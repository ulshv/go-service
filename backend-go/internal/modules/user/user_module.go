package user

type UserModule struct {
	UserService *UserService
}

func NewUserModule() *UserModule {
	userService := NewUserService()

	return &UserModule{
		UserService: userService,
	}
}
