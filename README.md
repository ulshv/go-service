# `go-service` is:

- Modular backend service written in Golang (mostly `stdlib`, but with some deps for SQL DBs).
- Includes User auth (register/login/authorization)
- Some simple domain-specific functionality as an example (implementation of online store or a marketplace)
- Can be started without external running deps (i.e. Postgres), SQLite can be used instead
- Headless backend, no GUI (maybe will be added in the future for as a reference implementation)
- [WIP] Full tests coverage (unit/e2e)
- [TODO] OpenAPI / Swagger

# Build & Run

Currently project doesn't have any UI client, so really the only way to test it's functionality is by automated tests. First, you need to install packages by `go get ./...` and then run tests (no PG instance running is required as it runs tests on SQLite DB): `go test ./...`. You can explore the DB state after the tests by exploring the `test.db` files in corresponding modules (i.e. `./internal/modules/auth/test.db`) via DBeaver or SQLite Viewer VSCode extension.

# A module code examples

App initialization:
```go
// `internal/application/application.go`

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
// `internal/modules/auth/auth_module.go`

func NewAuthModule(
  db *sqlx.DB,
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
// `internal/modules/auth/auth_service.go`

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
// `internal/modules/auth/auth_repository.go

func newAuthRepository(db *sqlx.DB) *authRepository {
  return &authRepository{
    db:     db,
    logger: logger.NewLogger("AuthRepository"),
  }
}
```

Module init in tests:
```go
// `internal/modules/auth/auth_test.go`

func initModule() *AuthModule {
  db := initDb()
  migrations.RunMigrations(db, "migrations", logger.NewLogger("Migrations"), database.SQLite)
  userModule := user.NewUserModule(db)
  authModule := NewAuthModule(userModule)
  return authModule
}

func TestRegister(t *testing.T) {
  module := initModule()
  ts := httptest.NewServer(http.HandlerFunc(module.AuthHandlers.registerHandler))
  defer ts.Close()

  tests := []struct{}
  // ...
}
```

Module handlers init/register:
```go
// `internal/modules/auth/auth_handlers.go`

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
// `internal/server/handlers.go`

func registerHandlers(mux *http.ServeMux, app *application.App) *http.ServeMux {
  app.AuthModule.AuthHandlers.RegisterHandlers(mux)
  app.UserModule.UserHandlers.RegisterHandlers(mux)
  // ...

  return mux
}
```
