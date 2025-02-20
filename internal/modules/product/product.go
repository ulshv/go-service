package product

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID        string          `json:"id" db:"id"`
	UserId    int             `json:"user_id" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Desc      string          `json:"desc" db:"desc"`
	Price     decimal.Decimal `json:"price" db:"price"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type createProductDto struct {
	Name  string          `json:"name"`
	Desc  string          `json:"desc"`
	Price decimal.Decimal `json:"price"`
}

func (p Product) MarshalJSON() ([]byte, error) {
	type Alias Product
	return json.Marshal(&struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		*Alias
	}{
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
		Alias:     (*Alias)(&p),
	})
}
