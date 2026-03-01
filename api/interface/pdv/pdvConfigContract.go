// Utilizado para gerenciar as configurações

package pdv

import (
	"context"
	u "myApi/helpers/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

type PDVConfig struct {
	Id                     int  `json:""`
	ConfirmToPinter        bool `json:"confirm_to_pinter"`
	BlockSaleNegativeStock bool `json:"block_sale_negative_stock"`
	ReserveStock           bool `json:"reserve_stock"`
}

func GetAll() ([]PDVConfig, error) {
	var configs []PDVConfig

	query := `
		SELECT
			id,
			confirm_to_pinter,
			block_sale_negative_stock,
			reserve_stock

		FROM
			config_pdv
	`

	rows, err := conn.Query(
		ctx,
		query,
	)

	if err != nil {
		u.GeneralLogger.Println("Erro ao executar a query:", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c PDVConfig

		if err := rows.Scan(
			&c.Id,
			&c.ConfirmToPinter,
			&c.BlockSaleNegativeStock,
			&c.ReserveStock,
		); err != nil {
			u.GeneralLogger.Println("Erro ao ler os dados da query:", err)
			return nil, err
		}

		configs = append(configs, c)
	}

	return configs, nil
}

func (c *PDVConfig) Update() (PDVConfig, error) {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transiction:", err)
		return PDVConfig{}, err
	}

	query := `
		UPDATE 	
			config_pdv

		SET
			confirm_to_pinter = $1,
			block_sale_negative_stock = $2,
			reserve_stock = $3
		
		WHERE
			id = 1
	`

	err = tx.QueryRow(
		ctx,
		query,
	).Scan(
		&c.Id,
		&c.ConfirmToPinter,
		&c.BlockSaleNegativeStock,
		&c.ReserveStock,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro executar o update:", err)
		return PDVConfig{}, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro no commit:", err)
		return PDVConfig{}, err
	}

	return *c, nil
}
