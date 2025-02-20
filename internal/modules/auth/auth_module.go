package auth

import "github.com/ulshv/go-service/internal/modules/user"

type AuthModule struct {
	service  *authService
	Handlers *authHandlers
}

func NewAuthModule(
	userModule *user.UserModule,
) *AuthModule {
	service := newAuthService(userModule.UserService)
	handlers := newAuthHandlers(service)

	return &AuthModule{
		service:  service,
		Handlers: handlers,
	}
}
