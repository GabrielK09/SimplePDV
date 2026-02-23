package customer

import (
	"context"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerContract struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CpfCnpj   string    `json:"cpf_cnpj"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (c CustomerContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	if c.Name == "" {
		errorsField["customer"] = "O nome do cliente precisa ser informado."
	}

	if len(c.CpfCnpj) > 14 {
		errorsField["cpf_cnpj"] = "O CPF/CNPJ do cliente não pode ser superior a 14 caracteres."
	}

	return errorsField
}

func GetAll() ([]CustomerContract, error) {
	var customers []CustomerContract

	query := `
		SELECT
			id,
			name,
			cpf_cnpj
			
		FROM
			customers
	`

	rows, err := conn.Query(
		ctx,
		query,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao fazer o select dos clientes: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c CustomerContract

		if err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.CpfCnpj,
		); err != nil {
			u.ErrorLogger.Println("Erro ao ler os dados do select dos clientes: ", err)
			return nil, err
		}

		customers = append(customers, c)
	}

	return customers, nil
}

func Show(id int) (*CustomerContract, error) {
	var c CustomerContract

	query := `
		SELECT
			id,
			Name,
			cpf_cnpj
			
		FROM
			customers

		WHERE
			id = $1
	`

	if err := conn.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&c.Id,
		&c.Name,
		&c.CpfCnpj,
	); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o select dos clientes: ", err)
		return nil, err
	}

	return &c, nil
}

func (c *CustomerContract) Create() error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transiction: ", err)
		return err
	}

	query := `
		INSERT INTO customers 
			(name, cpf_cnpj)

		VALUES
			($1, $2)

		RETURNING 
			id
	`

	if err = tx.QueryRow(
		ctx,
		query,
		&c.Name,
		&c.CpfCnpj,
	).Scan(&c.Id); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o insert: ", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o commit: ", err)
		return err
	}

	return nil
}
