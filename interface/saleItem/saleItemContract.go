package saleitem

import (
	"fmt"
	"log"
	"myApi/interface/product"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleItemContract []struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	Product   string    `json:"product"`
	Qtde      int       `json:"qtde"`
	SaleValue float64   `json:"sale_value"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SaleItem []SaleItemContract

var conn *pgxpool.Pool

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (i SaleItem) Validate() map[string]string {
	errorsField := make(map[string]string)

	for _, p := range i {
		if p. <= 0 {
			errorsField["sale_value"] = "O valor da venda não pode ser zerado."
		}

		if p.Qtde <= 0 {
			errorsField["qtde"] = "A quantidade do item da venda não pode ser zerado."
		}

		p, err := product.Show(p.ProductId)

		log.Println("P: ", p)

		if err != nil {
			errorsField["product_id"] = fmt.Sprintf("Um erro ocorreu: %s", err)
		}

		if p == nil {
			errorsField["product_id"] = "Produto não localizado."
		}

	}

	return errorsField
}
