package auth

import "github.com/ulshv/go-service/internal/modules/user"

type AuthModule struct {
	authService  *authService
	AuthHandlers *authHandlers
}

func NewAuthModule(
	userModule *user.UserModule,
) *AuthModule {
	service := newAuthService(userModule.UserService)
	handlers := newAuthHandlers(service)

	return &AuthModule{
		authService:  service,
		AuthHandlers: handlers,
	}
}
