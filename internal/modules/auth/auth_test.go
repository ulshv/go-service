package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ulshv/online-store-app/backend-go/internal/database"
	"github.com/ulshv/online-store-app/backend-go/internal/database/migrations"
	"github.com/ulshv/online-store-app/backend-go/internal/logger"
	"github.com/ulshv/online-store-app/backend-go/internal/modules/user"
	"github.com/ulshv/online-store-app/backend-go/internal/utils/testutils"
)

func initDb() *sqlx.DB {
	cfg := database.Config{
		Type:   database.SQLite,
		DBName: ":memory:",
	}
	db, err := database.NewConnection(cfg)
	if err != nil {
		panic(err)
	}
	return db
}

func initModule() *AuthModule {
	db := initDb()
	migrations.RunMigrations(db, "migrations", logger.NewLogger("Migrations"), database.SQLite)
	userModule := user.NewUserModule(db)
	return NewAuthModule(userModule)
}

func TestRegister(t *testing.T) {
	module := initModule()
	ts := httptest.NewServer(http.HandlerFunc(module.AuthController.registerHandler))
	defer ts.Close()

	tests := []struct {
		name       string
		payload    registerDto
		wantStatus int
		wantError  string
		wantResult registerResultDto
	}{
		{
			name: "registers new user",
			payload: registerDto{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusOK,
			wantError:  "",
			wantResult: registerResultDto{
				UserId: 1,
			},
		},
		{
			name: "don't register new user if email is already taken",
			payload: registerDto{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusConflict,
			wantError:  user.ErrEmailTaken.Error(),
			wantResult: registerResultDto{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}
			apiErr := testutils.ErrorStringFromBody(respBody)
			if apiErr != tt.wantError {
				t.Errorf("got error %s, want %s", apiErr, tt.wantError)
			}
			result := registerResultDto{}
			err = json.Unmarshal(respBody, &result)
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantResult.UserId != 0 && result.UserId != tt.wantResult.UserId {
				t.Errorf("got UserId %+v, want %+v", result, tt.wantResult)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	t.Skip()
}
