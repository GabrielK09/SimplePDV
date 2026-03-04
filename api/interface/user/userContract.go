package user

import (
	"context"
	"errors"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserContract struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Cpf       string    `json:"cpf"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool
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

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (u *UserContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	if u.Login == "" {
		errorsField["login"] = "Login ausente."
	}

	if u.Password == "" {
		errorsField["password"] = "Senha ausente."
	}

	return errorsField
}

func ShowByLogin(login string) (*UserContract, error) {
	var user UserContract

	query := `
		SELECT
			id,
			name,
			cpf,
			is_admin

		FROM
			users

		WHERE	
			login = $1
	`

	if err := conn.QueryRow(
		ctx,
		query,
		login,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Cpf,
		&user.IsAdmin,
	); err != nil {
		u.ErrorLogger.Println("Erro ao ler os dados da query: ", err)
		return nil, err
	}

	return &user, nil
}

func CreateDefaultUser() error {
	u.InfoLogger.Println("CreateDefaultUser started")

	var user UserContract

	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transiction: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	selectQuery := `
		SELECT
			id

		FROM
			users
			
		LIMIT 	
			1
	`

	if err = tx.QueryRow(
		ctx,
		selectQuery,
	).Scan(&user.Id); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao conferir se o usuário existe: ", err)
		return err
	}

	if user.Id > 0 {
		u.InfoLogger.Println("O usuário existe")
		return nil
	}

	u.GeneralLogger.Println("Não possui valores: ", user.Id)

	// Recolocar os dados do insert aqui

	if err != nil {
		return err
	}

	if _, err := tx.Exec(
		ctx,
		query,
		hashedPassword,
	); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o insert: ", err)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao fazer o commit: ", err)
		return err
	}

	u.GeneralLogger.Println("Usuário padrão cadastrado com sucesso!")

	return nil
}
