package saleitem

import (
	"fmt"
	"myApi/interface/product"
	"time"
)

type SaleItemContract []struct {
	Id        int       `json:"id"`
	SaleId    int       `json:"sale_id"`
	ProductId int       `json:"product_id"`
	Name      string    `json:"name"`
	Qtde      int       `json:"qtde"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (i SaleItemContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	for _, p := range i {
		if p.Price <= 0 {
			errorsField["sale_value"] = "O valor da venda não pode ser zerado."
		}

		if p.Qtde <= 0 {
			errorsField["qtde"] = "A quantidade do item da venda não pode ser zerado."
		}

		p, err := product.Show(p.ProductId)

		if err != nil {
			errorsField["product_id"] = fmt.Sprintf("Um erro ocorreu: %s", err)
		}

		if p == nil {
			errorsField["product_id"] = "Produto não localizado."
		}
	}

	return errorsField
}
