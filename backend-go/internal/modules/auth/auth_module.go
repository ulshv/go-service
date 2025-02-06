package auth

import "github.com/ulshv/online-store-app/backend-go/internal/modules/user"

type AuthModule struct {
	AuthService    *AuthService
	AuthController *AuthController
}

func NewAuthModule(
	userModule *user.UserModule,
) *AuthModule {
	service := newAuthService(userModule.UserService)
	controller := newAuthController(service)

	return &AuthModule{
		AuthService:    service,
		AuthController: controller,
	}
}
