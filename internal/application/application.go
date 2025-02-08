package application

import (
	"github.com/jmoiron/sqlx"
	"github.com/ulshv/online-store-app/backend-go/internal/database"
	"github.com/ulshv/online-store-app/backend-go/internal/modules/auth"
	"github.com/ulshv/online-store-app/backend-go/internal/modules/user"
)

type App struct {
	db         *sqlx.DB
	AuthModule *auth.AuthModule
	UserModule *user.UserModule
}

func NewApp(dbConfig database.Config) (*App, error) {
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		return nil, err
	}

	userModule := user.NewUserModule(db)
	authModule := auth.NewAuthModule(userModule)

	return &App{
		db:         db,
		AuthModule: authModule,
		UserModule: userModule,
	}, nil
}

func (a *App) Close() error {
	return a.db.Close()
}
