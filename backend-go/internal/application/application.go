package application

import (
	"github.com/ulshv/online-store-app/backend-go/internal/modules/auth"
	"github.com/ulshv/online-store-app/backend-go/internal/modules/user"
)

type App struct {
	AuthModule *auth.AuthModule
	UserModule *user.UserModule
}

func NewApp() *App {
	userModule := user.NewUserModule()
	authModule := auth.NewAuthModule(userModule)

	return &App{
		AuthModule: authModule,
		UserModule: userModule,
	}
}
