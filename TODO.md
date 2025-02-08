# go-service TODO

WIP/Done:
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
- [ ] Add auth/me route for getting user info based on the AuthToken
  - [x] /auth/me route and tests
  - [ ] Context WithValue
  - [ ] Middleware
  - [ ] Context get value and check it with token in the DB
- [ ] /auth/login handler (+ tests)

Backlog:
- [ ] Почистить проект и сделать код look prettier, возможно добавить побольше дебаг-логов
- [ ] Handle `LOG_LEVEL=debug` instead of `LOG_DEBUG=1`
- [ ] DB table auth_tokens: `id, token (unique), user_id` and revoke functionality for refresh_tokens
- [ ] After finishing auth/user modules, write down the architecture for the next features (like product/order)
- [ ] Tests for the user module

