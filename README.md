# `go-service` is:

- Modular backend service written in Golang (mostly `stdlib`, but with some deps for SQL DBs).
- Includes User auth (register/login/authorization)
- Some simple domain-specific functionality as an example (implementation of online store or a marketplace)
- Can be started without external running deps (i.e. Postgres), SQLite can be used instead
- Headless backend, no GUI (maybe will be added in the future for as a reference implementation)
- [WIP] Full tests coverage (unit/e2e)
- [TODO] OpenAPI / Swagger

# A module code examples

App initialization:
```go
db, err := database.NewConnection(dbConfig)
// handle err

userModule := user.NewUserModule(db)
authModule := auth.NewAuthModule(db, userModule)

return &App{
  db:         db,
  AuthModule: authModule,
  UserModule: userModule,
}, nil
```

Module init:
```go
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
```

Service init:
```go
func newAuthService(
  db *sqlx.DB,
  userService *user.UserService,
) *authService {
  return &authService{
    db:          *sqlx.DB,
    userService: userService,
    logger:      logger.NewLogger("AuthService"),
  }
}
```

Repository init:
```go
func newAuthRepository(db *sqlx.DB) *authRepository {
  return &authRepository{
    db:     db,
    logger: logger.NewLogger("AuthRepository"),
  }
}
```

Module handlers init/register:
```go
func newAuthHandlers(authService *authService) *authHandlers {
  return &authHandlers{
    authService: authService,
    logger:      logger.NewLogger("AuthHandlers"),
  }
}

func (h *authHandlers) RegisterHandlers(mux *http.ServeMux) {
  mux.HandleFunc("POST /api/v1/auth/register", h.registerHandler)
  mux.HandleFunc("POST /api/v1/auth/login", h.loginHandler)
}
```

Module handlers registration on `mux`:
```go
func registerHandlers(mux *http.ServeMux, app *application.App) *http.ServeMux {
  app.AuthModule.AuthHandlers.RegisterHandlers(mux)
  app.UserModule.UserHandlers.RegisterHandlers(mux)
  // ...

  return mux
}
```
