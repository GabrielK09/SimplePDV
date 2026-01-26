package product

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductContract struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Qtde      int       `json:"qtde"`
	Returned  int       `json:"returned"`
	Saled     int       `json:"saled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (p ProductContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	if p.Name == "" {
		errorsField["name"] = "O nome do produto é obrigatório!"
	}

	if p.Qtde <= 0 {
		errorsField["qtde"] = "A quantidade do produto é obrigatória!"
	}

	if p.Price <= 0 {
		errorsField["price"] = "O valor do produto é obrigatório!"
	}

	return errorsField
}

func (p *ProductContract) Create() error {
	query := `	
		INSERT INTO products
			(name, price, qtde, returned, saled)

		VALUES
			($1, $2, $3, $4, $5)
		
		RETURNING 
			id
	`

	err := conn.QueryRow(
		context.Background(),
		query,
		p.Name,
		p.Price,
		p.Qtde,
		p.Returned,
		p.Saled,
	).Scan(&p.Id)

	return err
}

func (p *ProductContract) Update() (ProductContract, error) {
	quey := `
		UPDATE
			products
		SET
			name = $2, 
			price = $3, 
			qtde = $4, 
			returned = $5, 
			saled = $6, 
			date_of_purchase = $7

		WHERE
			id = $1
		
		RETURNING
			id,
			name,
			price, 
			qtde, 
			returned, 
			saled, 
			date_of_purchase
	`

	err := conn.QueryRow(
		context.Background(),
		quey,
		p.Id,
		p.Name,
		p.Price,
		p.Qtde,
		p.Returned,
		p.Saled,
	).Scan(
		&p.Id,
		&p.Name,
		&p.Price,
		&p.Qtde,
		&p.Returned,
		&p.Saled,
	)

	if err != nil {
		return ProductContract{}, err
	}

	return *p, nil
}

func Show(id int) (*ProductContract, error) {
	query := `
		SELECT
			id,
			name, 
			price, 
			qtde, 
			returned, 
			saled, 
			date_of_purchase 
		FROM
			products

		WHERE
			id = $1
	`

	var p ProductContract

	err := conn.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&p.Id,
		&p.Name,
		&p.Price,
		&p.Qtde,
		&p.Returned,
		&p.Saled,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func GetAll() ([]ProductContract, error) {
	var products []ProductContract

	query := `
		SELECT
			id,
			name,
			price,
			qtde,           
			returned,
			saled,
			created_at,
			updated_at
		FROM
			products
	`

	rows, err := conn.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p ProductContract

		if err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Price,
			&p.Qtde,
			&p.Returned,
			&p.Saled,
		); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func Delete(id int) error {

	query := `
		DELETE FROM
			products

		WHERE id = $1

	`

	_, err := conn.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductContract) DiscountedQtde(qtde int) (ProductContract, error) {
	tx, err := conn.Begin(context.Background())

	if err != nil {
		return ProductContract{}, err
	}

	p, err = Show(p.Id)

	if err != nil {
		return ProductContract{}, err
	}

	query := `
		UPDATE	
			products p 

		SET
			p.qtde = p.qtde - $2,

		WHERE	
			p.id = $1
	`

	err = tx.QueryRow(
		context.Background(),
		query,
		p.Id,
		qtde,
	).Scan(
		&p.Id,
		&p.Qtde,
	)

	if err != nil {
		tx.Rollback(context.Background())
		return ProductContract{}, err
	}

	return *p, nil
}
