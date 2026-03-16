package dbaction

import (
	"context"
	"fmt"
	u "myApi/helpers/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBSeed(db *pgxpool.Pool, ctx context.Context) {
	fmt.Println("Executando seeder...")

	tx, err := db.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao iniciar a transiction Job.")
	}

	defer tx.Rollback(ctx)

	fmt.Println("Truncando tabela de produtos ...")
	if _, err := tx.Exec(
		ctx,
		`
			TRUNCATE TABLE products RESTART IDENTITY CASCADE;
		`,
	); err != nil {
		u.ErrorLogger.Fatal("Erro ao truncar tabela de produtos:", err)
	}

	query := `	
		INSERT INTO products
			(name, price, qtde, commission, returned, saled)

		VALUES
			(
				'Produto teste',
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
		query,
	); err != nil {
		u.ErrorLogger.Fatal("Erro ao criar o produto padrão:", err)
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Fatal("Erro ao commitar:", err)
	}

	fmt.Println("Produto padrão criado com sucesso!")
}
