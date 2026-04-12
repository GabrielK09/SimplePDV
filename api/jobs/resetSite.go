package jobs

import (
	"context"
	_ "embed"
	"fmt"
	u "myApi/helpers/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/resetSite.sql
var resetSql string

func ResetSite(db *pgxpool.Pool, ctx context.Context) {
	fmt.Println("Executando job ...")

	tx, err := db.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao iniciar a transiction Job.")
	}

	defer tx.Rollback(ctx)

	CreateUser(db, ctx)

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Fatal("Erro ao commitar o create do usuário.", err)
	}

	fmt.Println("Site resetado com sucesso!")
}
