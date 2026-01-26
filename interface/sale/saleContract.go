package sale

import (
	"context"
	"fmt"
	saleitem "myApi/interface/saleItem"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleContract struct {
	Id        int                         `json:"id"`
	Customer  string                      `json:"customer"`
	Specie    string                      `json:"specie"`
	SaleValue float64                     `json:"sale_value"`
	Status    string                      `json:"status"`
	Products  []saleitem.SaleItemContract `json:"products"`
	CreatedAt time.Time                   `json:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at"`
}

var conn *pgxpool.Pool

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (s SaleContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	if err := s.Products.Validate(); err != nil {
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

	querySale := `
		INSERT INTO sales
			(customer, specie, sale_value, status)

		VALUES
			($1, $2, $3, 'Pendente')

		RETURNING 
			id
	`

	err = tx.QueryRow(
		context.Background(),
		querySale,
		s.Customer,
		s.Specie,
		s.SaleValue,
	).Scan(&s.Id)

	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return err
}
