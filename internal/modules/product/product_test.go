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

	// first, register new User
	registerBody, err := json.Marshal(auth.RegisterDto{
		Email:    "test1@example.com",
		Password: "pass",
	})
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
	registerApiErr := testutils.ErrorCodeFromBody(registerRespBody)
	if registerApiErr != "" {
		t.Errorf("got register error '%s', want '%s'", registerApiErr, "")
	}
	registerResult := auth.RegisterResultDto{}
	err = json.Unmarshal(registerRespBody, &registerResult)
	if err != nil {
		t.Fatal(err)
	}
	createBody, err := json.Marshal(createProductDto{
		Name:  "first product",
		Desc:  "lorem ipsum something amenit...",
		Price: decimal.NewFromFloat(89.95),
	})
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
	if createResp.StatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", createResp.StatusCode, http.StatusOK)
	}
	createRespBody, err := io.ReadAll(createResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer createResp.Body.Close()
	apiErr := testutils.ErrorCodeFromBody(createRespBody)
	if apiErr != "" {
		t.Errorf("got error '%s', want '%s'", apiErr, "")
	}
	createResult := Product{}
	err = json.Unmarshal(createRespBody, &createResult)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("createResult:", createResult)
}
