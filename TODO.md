# go-service TODO

## Todo/Done:

### Initial
- [x] Make basic tests for /auth/register finally work
- [x] Reorganize the structure of the project, remove FE code (for now)
- [x] Switch to SQLite DB in user module and write some CRUD-operations in userRepository
- [x] Make /auth/register tests not to fail with the new SQLite-based userRepository
- [x] Update README
- [x] Rename the project to `go-service` (including GH repo) and adjust imports accordingly
- [x] Make the GH repo public
- [x] rename *controller to *handlers
- [x] Add code examples in the readme
- [x] Real generation of User.PasswordHash
- [x] Real generation of auth.AuthToken (with tests)
- [x] Add route where generated auth token can be used (in somewhere like `/product/create` or `/auth/me`)
- [x] Add auth/me route and tests for getting user info based on the AuthToken

### Auth
- [ ] Authenticating user's requests
  - [ ] Context WithValue (token or associated user)
  - [ ] Context GetValue (token or associated user)
  - [ ] Middleware to parse Authorization header and set associated user in the ctx

### User module
- [ ] /auth/login handler (+ tests)

### Product module
- [ ] add `user_id` for products
- [ ] use User from the ctx to set user_id or fail with unauthorized err response

### Logger
- [ ] Handle `LOG_LEVEL=debug` instead of `LOG_DEBUG=1`
- [ ] Move logger to pkg/logs package
- [ ] Instead of creating a new logger instance for every service, initialize it once and add WithService(svc string) method

## Backlog:
- [ ] Fix envutils/LoadEnvFiles (filenames is rewritten and arg is not used)
- [ ] Clean up the codebase and make the code look prettier
- [ ] DB table auth_tokens: `id, token (unique), user_id` and revoke functionality for refresh_tokens
- [ ] After finishing auth/user modules, write down the architecture for the next features (like product/order)
- [ ] Tests for the user module
- [ ] Add OpenAPI/Swagger
- [ ] Create Middleware chain handler
- [ ] (Maybe) add functionality so handlers don't directly write to the http.ResponseWriter,
  but rather return (any, error) and then add additional headers like (x-processing-time: X.XXms and handle some errors like pq/sqlite errors)
  and return their error messages only in dev environment, not in production (it will be just "Internal Server Error" on prod)
- [ ] Move bcrypt/JWT-related code to pkg/utils/authutils to abstract away them from the main app's logic

