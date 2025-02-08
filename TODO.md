# go-service

- [x] Доработать базовые тесты для /auth/register
- [x] Реорганизовать структуру проекта, убрать FE
- [x] Перейти на SQLite БД в модуле user и написать несколько CRUD-операциий в репозитории
- [x] Make /auth/register tests not to fail with the new SQLite-based userRepository
- [x] Обновить README
- [x] Переименовать проект в `go-service` (включая GH репу) и поправить импорты и GH репозиторий и поправить origin
- [x] Сделать GH репу публичной
- [x] rename *controller to *handlers
- [x] Add code examples in the readme
- [ ] Функционал /auth/login (+ тесты)
- [ ] Почистить проект и сделать код look prettier, возможно добавить побольше дебаг-логов
- [ ] Реальная генерация User.PasswordHash
- [ ] Реальная генерация auth.AuthToken
- [ ] Добавить роут где можно юзать сгенеренный auth token (типа `product/create`)
- [ ] Тесты для модуля user
- [ ] После завершения базового функционала модулей auth/user, Расписать архитектуру следующих модулей
