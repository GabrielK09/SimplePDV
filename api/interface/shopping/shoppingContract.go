package shopping

import (
	"context"
	"errors"
	"fmt"
	calchelper "myApi/helpers/calc"
	u "myApi/helpers/logger"
	"myApi/interface/product"
	productcharacteristics "myApi/interface/product/productCharacteristics"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShoppingContract struct {
	Id            int                    `json:"id"`
	Load          int                    `json:"load"`      // Carga - obrigatório
	Operation     string                 `json:"operation"` // Operação - opcional
	ShoppingItens []ShoppingItenContract `json:"shopping_itens"`
	TotalShopping float64                `json:"total_shopping"`
	Status        string                 `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type ShoppingItenContract struct {
	Id                               int                                                    `json:"id"`
	ShoppingId                       int                                                    `json:"shopping_id"`
	ProductId                        int                                                    `json:"product_id"`
	Name                             string                                                 `json:"name"`
	QtdePurchased                    int                                                    `json:"qtde_purchased"`
	PurchasedValue                   float64                                                `json:"purchased_value"`
	ShoppingItensWithCharacteristics *productcharacteristics.ProductCharacteristicsContract `json:"product_with_characteristics"`
	DeletedAt                        time.Time                                              `json:"deleted_at"`
	CreatedAt                        time.Time                                              `json:"created_at"`
	UpdatedAt                        time.Time                                              `json:"updated_at"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (s *ShoppingContract) Validate() map[string]string {
	errors := make(map[string]string)

	existingLoad, err := checkExistLoad(s.Load)

	if err != nil {
		errors["database"] = fmt.Sprintf("Erro ao conferir se a carga já existe: %s", err)
		return errors
	}

	if existingLoad {
		errors["load"] = "Essa carga já foi cadastrada."
		return errors
	}

	if s.Load < 0 {
		errors["load"] = "A informação da carga da compra é obrigatória."
		return errors
	}

	if s.TotalShopping < 0 {
		errors["total_shopping"] = "O valor total da compra não pode ser menor que zero."
		return errors
	}

	if len(s.ShoppingItens) <= 0 {
		errors["shopping_itens"] = "Itens ausentes na compra."
		return errors
	}

	for _, p := range s.ShoppingItens {
		product, err := product.Show(p.ProductId)

		if err != nil {
			errors["database"] = fmt.Sprintf("Erro ao conferir se o item existe: %s", err)
			return errors
		}

		if product == nil {
			errors["product"] = "Produto não localizado."
			return errors
		}

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
			total_shopping,
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
			&s.TotalShopping,
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
	var subTotal float64

	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transação: ", err)
		return 0, err
	}

	for _, p := range s.ShoppingItens {
		subTotal = calchelper.CalculateTotalSale(p.PurchasedValue, p.QtdePurchased)
	}

	if err := tx.QueryRow(
		ctx,
		`
			INSERT INTO shopping
				(
					load, 
					operation, 
					total_shopping
				)

			VALUES(
				$1, 
				'Entrada',
				$2
			)

			RETURNING
				id
		`,
		s.Load,
		subTotal,
	).Scan(
		&shoppingId,
	); err != nil {
		u.ErrorLogger.Println("Erro ao inserir os dados da compra: ", err)
		return 0, err
	}

	s.Id = shoppingId

	for _, p := range s.ShoppingItens {

		if p.ShoppingItensWithCharacteristics != nil {
			if _, err := tx.Exec(
				ctx,
				`
					INSERT INTO shopping_itens_grid
						(shopping_id, product_id, product_grid_id, size_saled, grid_qtde)
					VALUES
						($1, $2, $3, $4, $5)
				`,
				s.Id,
				p.ProductId,
				p.ShoppingItensWithCharacteristics.Id,
				p.ShoppingItensWithCharacteristics.Size,
				p.QtdePurchased,
			); err != nil {
				u.ErrorLogger.Println("Erro ao inserir os itens da compra: ", err)
				return 0, err
			}

			if _, err := tx.Exec(
				ctx,
				`
					INSERT INTO shopping_itens
						(shopping_id, product_id, name, qtde_purchased, purchased_value)

					VALUES($1, $2, $3, $4, $5)
				`,
				s.Id,
				p.ProductId,
				p.Name,
				p.QtdePurchased,
				p.PurchasedValue,
			); err != nil {
				u.ErrorLogger.Println("Erro ao inserir os itens da compra: ", err)
				return 0, err
			}

		}

		if p.ShoppingItensWithCharacteristics == nil {
			if _, err := tx.Exec(
				ctx,
				`
					INSERT INTO shopping_itens
						(shopping_id, product_id, name, qtde_purchased, purchased_value)

					VALUES($1, $2, $3, $4, $5)
				`,
				s.Id,
				p.ProductId,
				p.Name,
				p.QtdePurchased,
				p.PurchasedValue,
			); err != nil {
				u.ErrorLogger.Println("Erro ao inserir os itens da compra: ", err)
				return 0, err
			}

			/*
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							products

						SET
							name = $2,
							price = $3,
							qtde = qtde + $4

						WHERE
							id = $1
					`,
					p.ProductId,
					p.Name,
					p.PurchasedValue,
					p.QtdePurchased,
				); err != nil {
					u.ErrorLogger.Println("Erro ao alterar os itens no estoque, com base na compra: ", err)
					return 0, err
				}
			*/
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

	if err := conn.QueryRow(
		ctx,
		`
			SELECT
				id,
				load,
				operation,
				status,
				total_shopping
				
			FROM
				shopping

			WHERE
				id = $1
		`,
		shoppingId,
	).Scan(
		&s.Id,
		&s.Load,
		&s.Operation,
		&s.Status,
		&s.TotalShopping,
	); err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da compra: ", err)
		return nil, err
	}

	return &s, nil
}

func ShowShoppingItens(shopingId int) (*[]ShoppingItenContract, error) {
	var shoppingItens []ShoppingItenContract

	rows, err := conn.Query(
		ctx,
		`
			SELECT
				id,
				shopping_id,
				product_id,
				name,
				qtde_purchased,
				purchased_value
			FROM
				shopping_itens

			WHERE
				shopping_id = $1
		`,
		shopingId,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao executar a query: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item ShoppingItenContract

		if err := rows.Scan(
			&item.Id,
			&item.ShoppingId,
			&item.ProductId,
			&item.Name,
			&item.QtdePurchased,
			&item.PurchasedValue,
		); err != nil {
			u.ErrorLogger.Println("Erro ao executar a query: ", err)
			return nil, err
		}

		shoppingItens = append(shoppingItens, item)
	}

	return &shoppingItens, nil
}

func checkExistLoad(load int) (bool, error) {
	var shoppingId int

	err := conn.QueryRow(
		ctx,
		`
			SELECT
				id

			FROM
				shopping

			WHERE
				load = $1
		`,
		load,
	).Scan(&shoppingId)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao conferir se esse LOAD existe: ", err)
		return false, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	return false, nil
}

func ReturnLastShoppingLoad() (int, error) {
	var shoppingId int

	if err := conn.QueryRow(
		ctx,
		`
			SELECT
				load
			FROM
				shopping

			ORDER BY
				load DESC

			LIMIT 
				1
		`,
	).Scan(&shoppingId); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao conferir o ultimo ID do compras: ", err)
		return 0, err
	}

	return shoppingId, nil
}
