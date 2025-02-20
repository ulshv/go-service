package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ulshv/go-service/internal/core/database"
	"github.com/ulshv/go-service/internal/core/database/migrations"
	"github.com/ulshv/go-service/internal/modules/user"
	"github.com/ulshv/go-service/pkg/utils/testutils"
)

func initDb() *sqlx.DB {
	os.Remove("./test.db")
	cfg := database.Config{
		Type:   database.SQLite,
		DBName: "./test.db",
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		panic(err)
	}
	return db
}

func initModule() *AuthModule {
	db := initDb()
	migrations.RunMigrations(db, database.SQLite)
	userModule := user.NewUserModule(db)
	authModule := NewAuthModule(userModule)
	return authModule
}

func TestRegister(t *testing.T) {
	t.Skip()
	module := initModule()
	ts := httptest.NewServer(http.HandlerFunc(module.AuthHandlers.registerHandler))
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

func TestMe(t *testing.T) {
	module := initModule()
	ts1 := httptest.NewServer(http.HandlerFunc(module.AuthHandlers.registerHandler))
	ts2 := httptest.NewServer(http.HandlerFunc(module.AuthHandlers.meHandler))
	defer ts1.Close()
	defer ts2.Close()

	t.Run("creates new user and fetches /auth/me using Authorization header", func(t *testing.T) {
		body, err := json.Marshal(registerDto{
			Email:    "test@example.com",
			Password: "password123",
		})
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.Post(ts1.URL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		registerResult := registerResultDto{}
		err = json.NewDecoder(resp.Body).Decode(&registerResult)
		if err != nil {
			t.Fatal(err)
		}
		token := registerResult.Tokens.AccessToken
		req, err := http.NewRequest("GET", ts2.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
		}
		meData := user.User{}
		err = json.NewDecoder(resp.Body).Decode(&meData)
		if err != nil {
			t.Fatal(err)
		}
		if meData.Id != registerResult.UserId {
			t.Errorf("got UserId %d, want %d", meData.Id, registerResult.UserId)
		}
		if meData.Email != "test@example.com" {
			t.Errorf("got Email %s, want %s", meData.Email, "test@example.com")
		}
	})
}
