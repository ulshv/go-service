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

func initModule() *ProductModule {
	db := initDb()
	migrations.RunMigrations(db, database.SQLite)
	productModule := NewProductModule(db)
	return productModule
}

func TestCreateAndGetProduct(t *testing.T) {
	module := initModule()
	productSrv := httptest.NewServer(module.handlers.RegisterHandlers(&http.ServeMux{}))
	// getProductSrv := httptest.NewServer(http.HandlerFunc(module.handlers.createProductHandler))

	createBody, err := json.Marshal(createProductDto{
		Name:  "first product",
		Desc:  "lorem ipsum something amenit...",
		Price: decimal.NewFromFloat(89.95),
	})
	if err != nil {
		t.Fatal(err)
	}
	createUrl := fmt.Sprintf("%s/api/v1/products", productSrv.URL)
	createResp, err := http.Post(createUrl, "application/json", bytes.NewBuffer(createBody))
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
	apiErr := testutils.ErrorStringFromBody(createRespBody)
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
