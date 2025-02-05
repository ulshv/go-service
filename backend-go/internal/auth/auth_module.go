package auth

import "github.com/ulshv/online-store-app/backend-go/internal/user"

type AuthModule struct {
	authService *AuthService
}

func NewAuthModule(
	userModule *user.UserModule,
) *AuthModule {
	return &AuthModule{
		authService: NewAuthService(),
	}
}
