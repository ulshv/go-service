package product

import (
	"github.com/shopspring/decimal"
)

type productSvc struct {
	repo *productRepo
}

func newProductSvc(repo *productRepo) *productSvc {
	return &productSvc{
		repo: repo,
	}
}

func newProduct(name, desc string, price decimal.Decimal) Product {
	return Product{
		Name:  name,
		Desc:  desc,
		Price: price,
	}
}

func (s *productSvc) getProducts(offset, limit int) ([]Product, error) {
	return s.repo.list(offset, limit)
}

func (s *productSvc) getProductById(id int) (Product, error) {
	return s.repo.getById(id)
}

func (s *productSvc) createProduct(p Product) (Product, error) {
	return s.repo.create(p)
}

func (s *productSvc) updateProduct(p Product) error {
	return s.repo.update(p)
}
