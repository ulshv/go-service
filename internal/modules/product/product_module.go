package product

import "github.com/jmoiron/sqlx"

type ProductModule struct {
	handlers *productHandlers
}

func NewProductModule(db *sqlx.DB) *ProductModule {
	repo := newProductRepo(db)
	svc := newProductSvc(repo)
	handlers := newProductHandlers(svc)
	return &ProductModule{
		handlers: handlers,
	}
}
