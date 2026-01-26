package cashRegister

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CashRegisterContract struct {
	Id           int       `json:"id"`
	Description  string    `json:"description"`
	Customer     string    `json:"customer"`
	Specie       string    `json:"specie"`
	InputValue   float64   `json:"input_value"`
	OutputValue  float64   `json:"output_value"`
	TotalBalance float64   `json:"total_balance"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (c CashRegisterContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	if c.InputValue <= 0 && c.OutputValue <= 0 {
		errorsField["input_value"] = "O valor de entrada não pode ser menor que zero."
		errorsField["output_value"] = "O valor de saída não pode ser menor que zero."
	}

	if c.InputValue > 0 && c.OutputValue > 0 {
		errorsField["input_value"] = "O valor de entrada não pode ser informando quando houver um valor de saída informado."
		errorsField["output_value"] = "O valor de saída não pode ser informando quando houver um valor de entrada informado."
	}

	return errorsField
}

func GetLastId() (id int, err error) {
	query := `
		SELECT
			id
			
		FROM
			cash_registers

		ORDER BY 
			id DESC

	`

	err = conn.QueryRow(
		context.Background(),
		query,
	).Scan(
		&id,
	)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func getLastBalance() (float64, error) {
	query := `
		SELECT
			COALESCE(total_balance, 0)
			
		FROM
			cash_registers

		ORDER BY id DESC
		LIMIT 1
	`

	var balance float64

	err := conn.QueryRow(
		context.Background(),
		query,
	).Scan(
		&balance,
	)

	if balance == 0 {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return balance, nil
}

func GetAll() ([]CashRegisterContract, error) {
	var cashRegisters []CashRegisterContract

	query := `
		SELECT
			id,
			description,
			customer,
			specie,
			input_value,
			output_value,
			total_balance,
			created_at,
			updated_at 
		FROM
			cash_registers
	`

	rows, err := conn.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c CashRegisterContract

		if err := rows.Scan(
			&c.Id,
			&c.Description,
			&c.Customer,
			&c.Specie,
			&c.InputValue,
			&c.OutputValue,
			&c.TotalBalance,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		cashRegisters = append(cashRegisters, c)
	}

	return cashRegisters, nil
}

func (c *CashRegisterContract) Create(input_value, output_value float64) error {
	lastBalance, err := getLastBalance()

	if err != nil {
		return err
	}

	c.TotalBalance = lastBalance + c.InputValue - c.OutputValue

	query := `
		INSERT INTO cash_registers
			(description, customer, specie, input_value, output_value, total_balance)

		VALUES
			($1, $2, $3, $4, $5, $6)

		RETURNING id
	`

	err = conn.QueryRow(
		context.Background(),
		query,
		c.Description,
		c.Customer,
		c.Specie,
		c.InputValue,
		c.OutputValue,
		c.TotalBalance,
	).Scan(&c.Id)

	return err
}
