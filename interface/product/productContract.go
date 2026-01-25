package product

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool

type ProductContract struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Un             string    `json:"un"`
	Qtde           int       `json:"qtde"`
	Returned       int       `json:"returned"`
	Saled          int       `json:"saled"`
	DateOfPurchase time.Time `json:"date_of_purchase"`
}

func (p ProductContract) Validate() error {
	if p.Name == "" {
		return errors.New("o nome do produto é obrigatório!")
	}

	if p.Qtde <= 0 {
		return errors.New("a quantidade do produto é obrigatória!")
	}

	return nil
}

func (p ProductContract) Create() error {
	query := `	
		INSERT INTO products
			(name, un, qtde, returned, saled, date_of_purchase)

		VALUES
			($1, $2, $3, $4, $5, $6)
		
		RETURNING id
	`

	err := conn.QueryRow(
		context.Background(),
		query,
		p.Name,
		p.Un,
		p.Qtde,
		p.Returned,
		p.Saled,
		p.DateOfPurchase,
	).Scan(&p.Id)

	return err
}
