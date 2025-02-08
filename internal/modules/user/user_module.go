package user

import "github.com/jmoiron/sqlx"

type UserModule struct {
	UserService *UserService
}

func NewUserModule(db *sqlx.DB) *UserModule {
	userService := NewUserService(db)

	return &UserModule{
		UserService: userService,
	}
}
