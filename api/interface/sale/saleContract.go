package sale

import (
	"context"
	"fmt"
	"myApi/interface/product"
	saleitem "myApi/interface/saleItem"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleContract struct {
	Id        int                       `json:"id"`
	Customer  string                    `json:"customer"`
	Specie    string                    `json:"specie"`
	SaleValue float64                   `json:"sale_value"`
	Status    string                    `json:"status"`
	Products  saleitem.SaleItemContract `json:"products"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

var conn *pgxpool.Pool

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func calculateTotalSale(saleValue float64, qtde int) float64 {
	return saleValue * float64(qtde)
}

func (s SaleContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	if err := s.Products.Validate(); len(err) > 0 {
		errorsField["products"] = fmt.Sprintf("%s", err)
	}

	if s.SaleValue <= 0 {
		errorsField["sale_value"] = "O valor da venda não pode ser zerado."
	}

	return errorsField
}

func (s *SaleContract) Create() error {
	tx, err := conn.Begin(context.Background())

	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	querySale := `
		INSERT INTO sales
			(customer, specie, sale_value)

		VALUES
			($1, $2, $3)

		RETURNING 
			id,
			status
	`

	var saleId int

	err = tx.QueryRow(
		context.Background(),
		querySale,
		s.Customer,
		s.Specie,
		s.SaleValue,
	).Scan(
		&saleId,
		&s.Status,
	)

	s.Id = saleId

	if err != nil {
		return err
	}

	querySaleItem := `
		INSERT INTO sale_itens
			(product_id, name, qtde, sale_value, sale_id)

		VALUES
			($1, $2, $3, $4, $5)
			
		RETURNING 
			id,
			name,
			status
	`

	for idx := range s.Products {
		i := &s.Products[idx]

		totalSale := calculateTotalSale(i.SaleValue, i.Qtde)
		p, err := product.Show(i.ProductId)

		if err != nil {
			return err
		}

		err = tx.QueryRow(
			context.Background(),
			querySaleItem,
			i.ProductId,
			p.Name,
			i.Qtde,
			totalSale,
			saleId,
		).Scan(
			&i.Id,
			&i.Name,
			&i.Status,
		)

		i.Name = p.Name

		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return err
	}

	return nil
}
