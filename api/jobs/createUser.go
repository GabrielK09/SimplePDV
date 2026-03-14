package jobs

import (
	"context"
	_ "embed"
	"fmt"
	u "myApi/helpers/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

//go:embed sql/insert.sql
var insertSql string

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		14,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao criptografar a senha.", err)
		return "", err
	}

	return string(bytes), err
}

func CreateUser(db *pgxpool.Pool, ctx context.Context) {
	fmt.Println("Executando job ...")

	tx, err := db.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao iniciar a transiction Job.")
	}

	defer tx.Rollback(ctx)

	cryptedPass, err := hashPassword("123")

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao criptografar a senha.", err)
	}

	if _, err := tx.Exec(
		ctx,
		insertSql,
		cryptedPass,
	); err != nil {
		u.ErrorLogger.Fatal("Erro ao criar o usuário.", err)
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Fatal("Erro ao commitar o create do usuário.", err)
	}

	fmt.Println("Usuário criado com sucesso!")
}
