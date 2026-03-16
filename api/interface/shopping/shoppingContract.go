package shopping

import (
	"context"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShoppingContract struct {
	Id            int                    `json:"id"`
	Load          int                    `json:"load"`      // Carga - obrigatório
	Operation     string                 `json:"operation"` // Operação - opcional
	Status        string                 `json:"status"`
	ShoppingItens []ShoppingItenContract `json:"shopping_itens"`
	TotalShopping float64                `json:"total_shopping"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type ShoppingItenContract struct {
	Id             int       `json:"id"`
	ProductId      int       `json:"product_id"`
	Name           string    `json:"name"`
	QtdePurchased  int       `json:"qtde_purchased"`
	PurchasedValue float64   `json:"purchased_value"`
	ShoppingId     int       `json:"shopping_id"`
	Status         string    `json:"status"`
	DeletedAt      time.Time `json:"deleted_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (s *ShoppingContract) Validate() map[string]string {
	errors := make(map[string]string)

	if s.Load < 0 {
		errors["load"] = "A informação da carga da compra é obrigatória."
		return errors
	}

	if s.TotalShopping < 0 {
		errors["total_shopping"] = "O valor total da compra não pode ser menor que zero."
		return errors
	}

	for _, p := range s.ShoppingItens {
		if p.QtdePurchased <= 0 && p.PurchasedValue <= 0 {
			errors["qtde_purchased"] = "A qtde de compra não pode ser menor que zero."
			errors["purchased_value"] = "O valor do item da compra não pode ser menor que zero."
			return errors
		}
	}

	return errors
}

func GetAll() ([]ShoppingContract, error) {
	var shoppings []ShoppingContract

	querySelect := `
		SELECT
			id,
			load,
			operation,
			status
		FROM			
			shopping
	`

	rows, err := conn.Query(
		ctx,
		querySelect,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao ler os dados das compras: ", err)
		return []ShoppingContract{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var s ShoppingContract

		if err := rows.Scan(
			&s.Id,
			&s.Load,
			&s.Operation,
			&s.Status,
		); err != nil {
			u.ErrorLogger.Println("Erro ao ler os dados da query: ", err)
			return []ShoppingContract{}, err
		}

		shoppings = append(shoppings, s)
	}

	return shoppings, nil
}

func (s *ShoppingContract) Create() (int, error) {
	var shoppingId int
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transação: ", err)
		return 0, err
	}

	queryInsertShopping := `
		INSERT INTO shopping
			(load, operation)

		VALUES(
			$1, 
			'Entrada'
		)

		RETURNING
			id
	`

	if err := tx.QueryRow(
		ctx,
		queryInsertShopping,
		s.Load,
	).Scan(
		&shoppingId,
	); err != nil {
		u.ErrorLogger.Println("Erro ao inserir os dados da compra: ", err)
		return 0, err
	}

	s.Id = shoppingId

	queryInsertItens := `
		INSERT INTO 
			shopping_itens

		VALUES(
			shopping_id = $1,
			product_id = $2,
			name = $3,
			qtde_purchased = $4,
			purchased_value = $5
		)
	`

	queryForUpdateItens := `
		UPDATE
			products

		SET
			name = $2,
			price = $3,
			qtde = (
				SELECT
					COALESCE(SUM(qtde + $4), 0)

				FROM
					products

				WHERE
					id = $1
			),

		WHERE
			id = $1
	`

	for _, p := range s.ShoppingItens {
		if _, err := tx.Exec(
			ctx,
			queryInsertItens,
			shoppingId,
			p.ProductId,
			p.Name,
			p.QtdePurchased,
			p.PurchasedValue,
		); err != nil {
			u.ErrorLogger.Println("Erro ao inserir os itens da compra: ", err)
			return 0, err
		}

		if _, err := tx.Exec(
			ctx,
			queryForUpdateItens,
			p.ProductId,
			p.Name,
			p.PurchasedValue,
			p.QtdePurchased,
		); err != nil {
			u.ErrorLogger.Println("Erro ao alterar os itens no estoque, com base na compra: ", err)
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao salvar os dados da compra: ", err)
		return 0, err
	}

	return shoppingId, nil
}

func Show(shoppingId int) (*ShoppingContract, error) {
	var s ShoppingContract

	query := `
		SELECT
			id,
			load,
			operation,
			status,
			total_shopping,
			
		FROM
			shopping

		WHERE
			id = $1
	`

	if err := conn.QueryRow(
		ctx,
		query,
		shoppingId,
	).Scan(
		&s.Id,
		&s.Load,
		&s.Operation,
		&s.Status,
		&s.TotalShopping,
	); err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da venda")
		return nil, err
	}

	return &s, nil
}
