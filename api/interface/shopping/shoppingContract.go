package shopping

import (
	"context"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShoppingContract struct {
	Id            int                    `json:"id"`
	Load          int16                  `json:"load"`
	Operation     string                 `json:"operation"`
	Status        string                 `json:"status"`
	ShoppingItens []ShoppingItenContract `json:"shopping_itens"`
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

	if s.Status != "Pendente" {
		errors["status"] = "Essa compra já está finalizada."
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
	var shopping_id int
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transação: ", err)
		return 0, err
	}

	queryInsertShopping := `
		INSERT INTO 
			shopping
		VALUES(
			operation = 'Entrada'	
		)

		RETURNING
			id
	`

	if err := tx.QueryRow(
		ctx,
		queryInsertShopping,
	).Scan(
		&shopping_id,
	); err != nil {
		u.ErrorLogger.Println("Erro ao inserir os dados da compra: ", err)
		return 0, err
	}

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

	for _, p := range s.ShoppingItens {
		if _, err := tx.Exec(
			ctx,
			queryInsertItens,
			shopping_id,
			p.ProductId,
			p.Name,
			p.QtdePurchased,
			p.PurchasedValue,
		); err != nil {
			u.ErrorLogger.Println("Erro ao inserir os itens da ")
			return 0, err
		}
	}

	return 0, nil
}

func aggregateQtdeByName(newQtde int, productName string) error {

	return nil
}
