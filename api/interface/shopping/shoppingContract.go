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
	Id                               int                                                      `json:"id"`
	ShoppingId                       int                                                      `json:"shopping_id"`
	ProductId                        int                                                      `json:"product_id"`
	Name                             string                                                   `json:"name"`
	QtdePurchased                    int                                                      `json:"qtde_purchased"`
	PurchasedValue                   float64                                                  `json:"purchased_value"`
	ShoppingItensWithCharacteristics *[]productcharacteristics.ProductCharacteristicsContract `json:"product_with_characteristics"`
	DeletedAt                        time.Time                                                `json:"deleted_at"`
	CreatedAt                        time.Time                                                `json:"created_at"`
	UpdatedAt                        time.Time                                                `json:"updated_at"`
}

type ProductWithCharacteristics struct {
	Product         ShoppingItenContract                                    `json:"product"`
	Characteristics []productcharacteristics.ProductCharacteristicsContract `json:"characteristics"`
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

	defer tx.Rollback(ctx)

	for _, p := range s.ShoppingItens {
		subTotal = calchelper.CalculateTotalSale(p.PurchasedValue, p.QtdePurchased)
	}

	if err := tx.QueryRow(
		ctx,
		`
			INSERT INTO shopping
				(load, operation, total_shopping)

			VALUES($1, 'Entrada', $2)

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
			u.InfoLogger.Println("Produto possui caracteristicas, vai inserir na shopping_itens_grid e shopping_itens")

			var totalQtdeGrid int
			for _, grid := range *p.ShoppingItensWithCharacteristics {
				totalQtdeGrid += grid.GridQtde

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
					grid.Id,
					grid.Size,
					grid.GridQtde,
				); err != nil {
					u.ErrorLogger.Println("Erro ao inserir os itens da compra: ", err)
					return 0, err
				}
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
				totalQtdeGrid,
				p.PurchasedValue,
			); err != nil {
				u.ErrorLogger.Println("Erro ao inserir os itens da compra: ", err)
				return 0, err
			}
		}

		if p.ShoppingItensWithCharacteristics == nil {
			u.InfoLogger.Println("Produto não possui caracteristicas, vai inserir na shopping_itens")

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
	); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao pegar os dados da compra: ", err)
		return nil, err
	}

	return &s, nil
}

func ShowShoppingItens(shopingId int) (*[]ShoppingItenContract, error) {
	var shoppingItens []ShoppingItenContract

	rowsShoppingItens, err := conn.Query(
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

	defer rowsShoppingItens.Close()

	for rowsShoppingItens.Next() {
		var item ShoppingItenContract

		if err := rowsShoppingItens.Scan(
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

func ShowShoppingGridItens(shopingId, productId int) (*[]productcharacteristics.ProductCharacteristicsContract, error) {
	var shoppingItens []productcharacteristics.ProductCharacteristicsContract

	rowsShoppingGridItens, err := conn.Query(
		ctx,
		`
			SELECT
				product_id,
				size_saled,
				grid_qtde
			FROM
				shopping_itens_grid

			WHERE
				shopping_id = $1 AND
				product_id = $2
		`,
		shopingId,
		productId,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao executar a query: ", err)
		return nil, err
	}

	defer rowsShoppingGridItens.Close()

	for rowsShoppingGridItens.Next() {
		var item productcharacteristics.ProductCharacteristicsContract

		if err := rowsShoppingGridItens.Scan(
			&item.ProductId,
			&item.Size,
			&item.GridQtde,
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

	return true, nil
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

func (s *ShoppingContract) UpdateShopping() error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transação: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	for _, p := range s.ShoppingItens {
		if p.ShoppingItensWithCharacteristics != nil {
			u.InfoLogger.Println("O produto possui caracteristicas, vai alterar a shopping_itens_grid e shopping_itens")

			for _, c := range *p.ShoppingItensWithCharacteristics {
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							shopping_itens_grid
					
						SET
							grid_qtde = $2
		
						WHERE
							shopping_id = $1 AND 
							size_saled = $3 AND
							product_id = $4
							
					`,
					s.Id,
					c.GridQtde,
					c.Size,
					p.ProductId,
				); err != nil {
					u.ErrorLogger.Println("Erro ao alterar os dados da compra:", err)
					return err
				}
			}

			for _, c := range *p.ShoppingItensWithCharacteristics {
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							shopping_itens
					
						SET
							qtde_purchased = (
								SELECT
									COALESCE(SUM(grid_qtde), 0)
	
								FROM
									shopping_itens_grid
	
								WHERE
									shopping_id = $1 AND 
									size_saled = $2 AND 
									product_id = $3
							)
								
						WHERE
							shopping_id = $1 AND 
							product_id = $3
					`,
					s.Id,
					c.Size,
					p.ProductId,
				); err != nil {
					u.ErrorLogger.Println("Erro ao alterar os dados da compra:", err)
					return err
				}
			}
		}

		if p.ShoppingItensWithCharacteristics == nil {
			u.InfoLogger.Println("O produto não possui caracteristicas, vai alterar shopping_itens")

			if _, err := tx.Exec(
				ctx,
				`
					UPDATE
						shopping_itens
				
					SET
						qtde_purchased = $2,
						Purchased_Value = $3
	
					WHERE
						shopping_id = $1 AND
						product_id = $4
						
				`,
				s.Id,
				p.QtdePurchased,
				p.PurchasedValue,
				p.ProductId,
			); err != nil {
				u.ErrorLogger.Println("Erro ao alterar os dados da compra:", err)
				return err
			}
		}
	}

	u.InfoLogger.Println("Vai alterar o valor total da compra")
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE 
				shopping
			SET
				total_shopping = (
					SELECT
						COALESCE(SUM(qtde_purchased * purchased_value), 0)

					FROM
						shopping_itens

					WHERE
						shopping_id = $1
				)
			WHERE
				id = $1
		`,
		s.Id,
	); err != nil {
		u.ErrorLogger.Println("Erro ao alterar o total da compra depois da inserção/alteração dos itens: ", err)
		return err
	}

	u.InfoLogger.Println("Valor da compra atualizado.")

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao comitar: ", err)

		return err
	}

	return nil
}
