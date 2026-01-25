package saleitem

import "time"

type SaleItemContract struct {
	Id             int       `json:"id"`
	ProductId      int       `json:"product_id"`
	Product        string    `json:"product"`
	Qtde           int       `json:"qtde"`
	SaleValue      float64   `json:"sale_value"`
	DateOfMovement time.Time `json:"date_of_movement"`
}
