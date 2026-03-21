package dbaction

import (
	"context"
	"errors"
	"fmt"
	u "myApi/helpers/logger"
	"myApi/interface/customer"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DBSeed(db *pgxpool.Pool, ctx context.Context) {
	fmt.Println("Executando seeder...")

	tx, err := db.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao iniciar a transiction Job.")
	}

	defer tx.Rollback(ctx)

	var c customer.CustomerContract

	fmt.Println("Truncando tabela de produtos e clientes ...")
	if _, err := tx.Exec(
		ctx,
		`
			TRUNCATE TABLE shopping, products, customers RESTART IDENTITY CASCADE;
		`,
	); err != nil {
		u.ErrorLogger.Fatal("Erro ao truncar tabela de produtos e ou clientes:", err)
	}

	queryForInsertProducts := `	
		INSERT INTO products
			(name, price, qtde, commission, returned, saled)

		VALUES
			(
				'PRODUTO TESTE',
				10,
				10,
				0,
				0,
				0
			)
		
		RETURNING 
			id
	`

	if _, err := tx.Exec(
		ctx,
		queryForInsertProducts,
	); err != nil {
		u.ErrorLogger.Fatal("Erro ao criar o produto padrão:", err)
	}

	selectForVerifyExistsCustomerQuery := `
		SELECT
			id

		FROM
			customers

		LIMIT 1
	`

	if err := tx.QueryRow(
		ctx,
		selectForVerifyExistsCustomerQuery,
	).Scan(&c.Id); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Fatal("Erro ao conferir se o cliente padrão existe: ", err)

	}

	if c.Id > 0 {
		u.InfoLogger.Fatal("O cliente existe")
	}

	queryForInsertCustomer := `
		INSERT INTO customers 
			(id, name, cpf_cnpj)
		VALUES
			(1, 'CONSUMIDOR PADRÃO', '')
	`

	if _, err := tx.Exec(
		ctx,
		queryForInsertCustomer,
	); err != nil {
		u.ErrorLogger.Fatal("Erro ao fazer o insert: ", err)
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Fatal("Erro ao commitar:", err)
	}

	fmt.Println("Produto padrão criado com sucesso!")
}
