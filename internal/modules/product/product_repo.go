package product

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type productRepo struct {
	db *sqlx.DB
}

var queries = struct {
	getProducts    string
	getProductById string
	createProduct  string
	updateProduct  string
}{
	getProducts:    "SELECT * from products LIMIT :limit OFFSET :offset",
	getProductById: "SELECT * FROM products WHERE id = $1",
	createProduct: `INSERT INTO products (name, desc, price, created_at, updated_at)
VALUES (:name, :desc, :price, :created_at, :updated_at)
RETURNING id, name, desc, price, created_at, updated_at`,
	updateProduct: `UPDATE products
SET name = :name, desc = :desc, price = :price, updated_at = :updated_at
WHERE id = :id
RETURNING id, name, desc, price, created_at, updated_at`,
}

func newProductRepo(db *sqlx.DB) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) list(offset, limit int) ([]Product, error) {
	var p []Product
	err := r.db.Select(&p, queries.getProducts)
	if err != nil {
		return []Product{}, err
	}
	return p, nil
}

func (r *productRepo) getById(id int) (Product, error) {
	var p Product
	err := r.db.Get(&p, queries.getProductById, id)
	if err != nil {
		return Product{}, err
	}
	return p, nil
}

func (r *productRepo) create(p Product) (Product, error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	rows, err := r.db.NamedQuery(queries.createProduct, p)
	if err != nil {
		return Product{}, err
	}
	defer rows.Close()
	var created Product
	rows.Next()
	err = rows.StructScan(&created)
	if err != nil {
		return Product{}, err
	}
	return created, nil
}

func (r *productRepo) update(p Product) error {
	p.UpdatedAt = time.Now()
	_, err := r.db.NamedExec(queries.updateProduct, p)
	return err
}
