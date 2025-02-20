package product

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
	"github.com/shopspring/decimal"
	"github.com/ulshv/go-service/internal/core/database"
	"github.com/ulshv/go-service/internal/core/database/migrations"
	"github.com/ulshv/go-service/internal/core/httperrs"
	"github.com/ulshv/go-service/internal/modules/auth"
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

type initModulesResult struct {
	productModule *ProductModule
	authModule    *auth.AuthModule
}

func initModules() initModulesResult {
	db := initDb()
	migrations.RunMigrations(db, database.SQLite)
	userModule := user.NewUserModule(db)
	authModule := auth.NewAuthModule(userModule)
	productModule := NewProductModule(db)
	return initModulesResult{
		authModule:    authModule,
		productModule: productModule,
	}
}

func TestCreateAndGetProduct(t *testing.T) {
	modules := initModules()
	mux := &http.ServeMux{}
	modules.authModule.Handlers.RegisterHandlers(mux)
	modules.productModule.handlers.RegisterHandlers(mux)
	srv := httptest.NewServer(mux)

	t.Run("basic register / create product", func(t *testing.T) {
		tests := []struct {
			name                    string
			registerDto             auth.RegisterDto
			registerApiErrCode      string
			registerUserId          int
			createProductDto        createProductDto
			createProductApiErrCode string
		}{
			{
				name: "create a new user and create a product",
				registerDto: auth.RegisterDto{
					Email:    "test1@example.com",
					Password: "pass",
				},
				registerApiErrCode: "",
				registerUserId:     1,
				createProductDto: createProductDto{
					Name:  "first product",
					Desc:  "lorem ipsum something amenit...",
					Price: decimal.NewFromFloat(89.95),
				},
			},
			{
				name: "create another new user and create a product for this user",
				registerDto: auth.RegisterDto{
					Email:    "test2@example.com",
					Password: "pass",
				},
				registerApiErrCode: "",
				registerUserId:     2,
				createProductDto: createProductDto{
					Name:  "secodn product",
					Desc:  "just ipsum amenit",
					Price: decimal.NewFromFloat(19.45),
				},
			},
			{
				name: "try to create a new user with already taken email and create a product",
				registerDto: auth.RegisterDto{
					Email:    "test2@example.com",
					Password: "pass",
				},
				registerApiErrCode: httperrs.ErrCodeEmailTaken,
				registerUserId:     0,
				createProductDto: createProductDto{
					Name:  "secodn product",
					Desc:  "just ipsum amenit",
					Price: decimal.NewFromFloat(19.45),
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// first, register new User
				registerBody, err := json.Marshal(tt.registerDto)
				if err != nil {
					t.Fatal(err)
				}
				registerUrl := fmt.Sprintf("%s/api/v1/auth/register", srv.URL)
				registerResp, err := http.Post(registerUrl, "application/json", bytes.NewBuffer(registerBody))
				if err != nil {
					t.Fatal(err)
				}
				registerRespBody, err := io.ReadAll(registerResp.Body)
				if err != nil {
					t.Fatal(err)
				}
				defer registerResp.Body.Close()
				registerApiErrCode := testutils.ErrorCodeFromBody(registerRespBody)
				fmt.Println("registerApiErrCode", registerApiErrCode)
				if tt.registerApiErrCode != "" {
					if registerApiErrCode != tt.registerApiErrCode {
						t.Fatalf("got register error code '%s', want '%s'", registerApiErrCode, tt.registerApiErrCode)
					}
					return
				}
				registerResult := auth.RegisterResultDto{}
				err = json.Unmarshal(registerRespBody, &registerResult)
				if err != nil {
					t.Fatal(err)
				}
				if registerResult.UserId != tt.registerUserId {
					t.Fatalf("got register userId %v, want %v", registerResult.UserId, tt.registerUserId)
				}
				createBody, err := json.Marshal(tt.createProductDto)
				if err != nil {
					t.Fatal(err)
				}
				createUrl := fmt.Sprintf("%s/api/v1/products", srv.URL)
				req, err := http.NewRequest("POST", createUrl, bytes.NewBuffer(createBody))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", registerResult.Tokens.AccessToken))
				createResp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}
				// if createResp.StatusCode != http.StatusOK {
				// 	t.Errorf("got status %d, want %d", createResp.StatusCode, http.StatusOK)
				// }
				createRespBody, err := io.ReadAll(createResp.Body)
				if err != nil {
					t.Fatal(err)
				}
				defer createResp.Body.Close()
				createProductApiErrCode := testutils.ErrorCodeFromBody(createRespBody)
				if createProductApiErrCode != tt.createProductApiErrCode {
					t.Fatalf("got create product api error code '%s', want '%s'", createProductApiErrCode, tt.createProductApiErrCode)
				}
				createResult := Product{}
				err = json.Unmarshal(createRespBody, &createResult)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println("createResult:", createResult)
				if createResult.UserId != registerResult.UserId {
					t.Fatalf("got userId %v, want %v", createResult.UserId, registerResult.UserId)
				}
				if createResult.Name != tt.createProductDto.Name ||
					createResult.Desc != tt.createProductDto.Desc ||
					!createResult.Price.Equal(tt.createProductDto.Price) {
					t.Fatalf("got create product result '%+v', want '%+v'", createResult, tt.createProductDto)
				}
			})
		}
	})
}
