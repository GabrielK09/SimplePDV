package cashRegister

import (
	"context"
	"fmt"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CashRegisterContract struct {
	Id           int       `json:"id"`
	Description  string    `json:"description"`
	CustomerId   int       `json:"customer_id"`
	Customer     string    `json:"customer"`
	SpecieId     int       `json:"specie_id"`
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

func getLastBalance(tx pgx.Tx) (float64, error) {
	query := `
		SELECT
			COALESCE(total_balance, 0)
			
		FROM
			cash_registers

		ORDER BY 
			id DESC
		LIMIT 1
	`

	var balance float64

	err := tx.QueryRow(
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
			customer_id,
			customer,
			specie_id,
			specie,
			input_value,
			output_value,
			total_balance
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
			&c.CustomerId,
			&c.Customer,
			&c.SpecieId,
			&c.Specie,
			&c.InputValue,
			&c.OutputValue,
			&c.TotalBalance,
		); err != nil {
			return nil, err
		}

		cashRegisters = append(cashRegisters, c)
	}

	return cashRegisters, nil
}

func (c *CashRegisterContract) Create(
	tx pgx.Tx,
	input_value,
	output_value float64,
	shoppingId,
	saleId,
	cusotmerId int,
	customer string,
) map[string]string {
	errorsField := make(map[string]string)

	c.InputValue = input_value
	c.OutputValue = output_value
	c.CustomerId = cusotmerId
	c.Customer = customer

	if c.InputValue > 0 && c.OutputValue > 0 {
		errorsField["input_value"] = "Um registro no caixa não pode ter um valor de entrada no mesmo registro de uma saída."
		errorsField["output_value"] = "Um registro no caixa não pode ter um valor de saída no mesmo registro de uma entrada."
	}

	if shoppingId > 0 && saleId > 0 {
		errorsField["shopping_id"] = "Uma venda e uma compra não podem ser gravadas no mesmo registro do caixa."
	}

	lastBalance, err := getLastBalance(tx)

	if err != nil {
		u.ErrorLogger.Println("Erro no getLastBalance (cash-register-contract): ", err)
		errorsField["error"] = fmt.Sprintf("%s", err)
	}

	u.GeneralLogger.Printf("Dados: lastBalance - %f|InputValue - %f|OutputValue - %f", lastBalance, c.InputValue, c.OutputValue)

	c.TotalBalance = lastBalance + c.InputValue - c.OutputValue

	var sale interface{} = nil
	var shopping interface{} = nil
	var description interface{} = ""

	if saleId > 0 {
		sale = saleId

		if input_value > 0 {
			description = fmt.Sprintf("Venda n° %d", saleId)
		} else {
			description = fmt.Sprintf("Estorno de venda n° %d", saleId)
		}
	}

	if shoppingId > 0 {
		shopping = shoppingId

		if input_value > 0 {
			description = fmt.Sprintf("Compra n° %d", saleId)
		} else {
			description = fmt.Sprintf("Estorno de compra n° %d", saleId)
		}
	}

	if shoppingId <= 0 && saleId <= 0 {
		description = "Registro manual do caixa"
	}

	if len(errorsField) > 0 {
		return errorsField
	}

	args := []interface{}{
		description,
		c.CustomerId,
		c.Customer,
		c.SpecieId,
		c.Specie,
		c.InputValue,
		c.OutputValue,
		c.TotalBalance,
		sale,
		shopping,
	}

	u.GeneralLogger.Println("Dados do insert: ", args)

	query := `
		INSERT INTO cash_registers
			(
				description, 
				customer_id, 
				customer, 
				specie_id, 
				specie, 
				input_value, 
				output_value, 
				total_balance,
				sale_id,
				shopping_id
			)

		VALUES
			(
				$1, 
				$2, 
				$3, 
				$4, 
				$5, 
				$6,
				$7,
				$8,
				$9,
				$10
			)

		RETURNING id
	`

	err = tx.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(&c.Id)

	if err != nil {
		u.ErrorLogger.Println("Erro no create (cash-register-contract): ", err)
		errorsField["database"] = err.Error()
		return errorsField
	}

	return errorsField
}
