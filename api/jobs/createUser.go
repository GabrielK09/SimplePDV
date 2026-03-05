package jobs

import (
	"context"
	_ "embed"
	"fmt"
	"myApi/db"
	u "myApi/helpers/logger"

	"golang.org/x/crypto/bcrypt"
)

//go:embed sql/insert.sql
var insertSql string

var ctx = context.Background()

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

func CreateUser() {
	fmt.Println("Executando job ...")
	db, err := db.Init()

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao iniciar o banco de dados no Job.")
	}

	defer db.Close()

	tx, err := db.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao iniciar a transiction Job.")
	}

	defer tx.Rollback(ctx)

	cryptedPass, err := hashPassword("123")

	if err != nil {
		u.ErrorLogger.Fatal("Erro ao criptografar a senha.", err)
	}

	if _, err := db.Exec(
		ctx,
		insertSql,
		cryptedPass,
	); err != nil {
		u.ErrorLogger.Fatal("Erro ao criar o usuário.", err)
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Fatal("Erro ao commitar o create do usuário.", err)
	}

	db.Close()

	fmt.Println("Usuário criado com sucesso!")
}
