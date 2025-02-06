package auth

import "github.com/ulshv/online-store-app/backend-go/internal/modules/user"

type AuthModule struct {
	authService    *authService
	AuthController *authController
}

func NewAuthModule(
	userModule *user.UserModule,
) *AuthModule {
	service := newAuthService(userModule.UserService)
	controller := newAuthController(service)

	return &AuthModule{
		authService:    service,
		AuthController: controller,
	}
}
