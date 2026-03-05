package user

import (
	"context"
	"errors"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
			login,
			password,
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
		&user.Login,
		&user.Password,
		&user.IsAdmin,
	); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao ler os dados da query: ", err)
		return nil, err
	}

	return &user, nil
}
