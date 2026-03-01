package customer

import (
	"context"
	"errors"
	"fmt"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5"
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

func Delete(id int) error {
	verify := `
		SELECT
			id
		FROM
			sale_itens
		WHERE
			customer_id = $1
		LIMIT
			1
	`

	var saleCustomerId int

	err := conn.QueryRow(
		ctx,
		verify,
		saleCustomerId,
	).Scan(
		&saleCustomerId,
	)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			u.ErrorLogger.Println("Erro ao conferir se possui cadastro em vendas: ", err)
			return err
		}

	} else {
		return fmt.Errorf("Cliente já cadastradao em uma venda.")
	}

	query := `
		DELETE FROM
			customers

		WHERE 
			id = $1
	`

	_, err = conn.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao deletar: ", err)
		return err
	}

	return nil
}

func (c *CustomerContract) Update() (CustomerContract, error) {
	quey := `
		UPDATE
			customers
		SET
			name = $2, 
			cpf_cnpj = $3

		WHERE
			id = $1
		
		RETURNING
			id,
			name,
			cpf_cnpj
	`

	err := conn.QueryRow(
		context.Background(),
		quey,
		c.Id,
		c.Name,
		c.CpfCnpj,
	).Scan(
		&c.Id,
		&c.Name,
		&c.CpfCnpj,
	)

	if err != nil {
		return CustomerContract{}, err
	}

	return *c, nil
}

func CreateDefaultCustomer() error {
	u.InfoLogger.Println("CreateDefaultCustomer started")

	var c CustomerContract
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transiction: ", err)
		return err
	}

	selectQuery := `
		SELECT
			id

		FROM
			customers

		LIMIT 1
	`

	if err := tx.QueryRow(
		ctx,
		selectQuery,
	).Scan(&c.Id); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao conferir se o cliente padrão existe: ", err)
		return err
	}

	if c.Id > 0 {
		u.InfoLogger.Println("O cliente existe")
		return nil
	}

	u.GeneralLogger.Println("Não possui valores: ", c.Id)

	query := `
		INSERT INTO customers 
			(id, name, cpf_cnpj)
		VALUES
			(1, 'Consumidor padrão', '')
	`

	if _, err = tx.Exec(
		ctx,
		query,
	); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o insert: ", err)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao fazer o commit: ", err)
		return err
	}

	u.GeneralLogger.Println("Cliente padrão cadastrado com sucesso!")

	return nil
}
