package product

import (
	"context"
	"errors"
	"fmt"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductContract struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Qtde       int       `json:"qtde"`
	Commission float64   `json:"commission"`
	Returned   int       `json:"returned"`
	Saled      int       `json:"saled"`
	UseGrid    bool      `json:"use_grid"`
	DeletedAt  time.Time `json:"deleted_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

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

	if p.Commission < 0 {
		errorsField["commission_min"] = "O valor de comissão não pode ser menor que zero."
	}

	if p.Commission > 100 {
		errorsField["commission_max"] = "O valor de comissão não pode ser maior que 100%."
	}

	return errorsField
}

func (p *ProductContract) Create() error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transação: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	query := `	
		INSERT INTO products
			(name, price, qtde, commission, use_grid)

		VALUES
			($1, $2, $3, $4, $5)
	`

	if _, err := tx.Exec(
		ctx,
		query,
		p.Name,
		p.Price,
		p.Qtde,
		p.Commission,
		p.UseGrid,
	); err != nil {
		u.ErrorLogger.Println("Erro ao inserir o novo produto: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao commitar: ", err)
		return err
	}

	return err
}

func (p *ProductContract) Update() error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transição: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				products
			SET
				name = $2, 
				price = $3, 
				qtde = $4, 
				commission = $5,
				use_grid = $6

			WHERE
				id = $1
		`,
		p.Id,
		p.Name,
		p.Price,
		p.Qtde,
		p.Commission,
		p.UseGrid,
	); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o update do produto: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao commitar: ", err)
		return err
	}

	return nil
}

func Show(id int) (*ProductContract, error) {
	query := `
		SELECT
			id,
			name, 
			price, 
			qtde, 
			commission, 
			use_grid 
		FROM
			products

		WHERE
			id = $1
	`

	var p ProductContract

	if err := conn.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&p.Id,
		&p.Name,
		&p.Price,
		&p.Qtde,
		&p.Commission,
		&p.UseGrid,
	); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &p, nil
}

func ShowByName(productName string) ([]ProductContract, error) {
	var products []ProductContract

	u.InfoLogger.Println("ShowByName:", productName)

	query := `
		SELECT
			id,
			name, 
			price, 
			qtde, 
			commission
		FROM
			products

		WHERE
			name ILIKE '%'||$1||'%'

		ORDER BY 
			name

		LIMIT 
			20
	`

	rows, err := conn.Query(
		context.Background(),
		query,
		productName,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao realizar a query:", err)
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
			&p.Commission,
		); err != nil {
			u.ErrorLogger.Println("Erro ao realizar a query:", err)
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func GetAll() ([]ProductContract, error) {
	var products []ProductContract

	query := `
		SELECT
			id,
			name,
			price,
			qtde,           
			commission,
			use_grid
		FROM
			products

		WHERE
			deleted_at IS NULL
	`

	rows, err := conn.Query(
		context.Background(),
		query,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro: ", err)
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
			&p.Commission,
			&p.UseGrid,
		); err != nil {
			u.ErrorLogger.Println("Erro: ", err)
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func Delete(id int, deletedAt time.Time) error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transição: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	product, err := Show(id)

	if err != nil {
		u.ErrorLogger.Println("Ocorreu um erro ao consultar o produto: ", err)
		return err
	}

	if product == nil {
		u.ErrorLogger.Println("Produto não localizado.")
		return fmt.Errorf("Produto não localizado.")
	}

	queryUpdateDeletedAt := `
		UPDATE
			products

		SET
			deleted_at = $2

		WHERE 
			id = $1
	`

	if _, err = conn.Exec(
		ctx,
		queryUpdateDeletedAt,
		id,
		deletedAt,
	); err != nil {
		u.ErrorLogger.Println("Erro ao deletar o produto: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao commitar: ", err)
		return err
	}

	return nil
}

func (p *ProductContract) DiscountedQtde(ctx context.Context, tx pgx.Tx, qtde int) error {
	if qtde <= 0 {
		return fmt.Errorf("Qtde inválida: %d", qtde)
	}

	query := `
		UPDATE	
			products 

		SET
			qtde = qtde - $2

		WHERE	
			id = $1

		RETURNING
			id,
			qtde,
			name
	`

	return tx.QueryRow(ctx, query, p.Id, qtde).Scan(&p.Id, &p.Qtde, &p.Name)
}
