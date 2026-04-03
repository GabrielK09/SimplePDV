package sale

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	calchelper "myApi/helpers/calc"
	u "myApi/helpers/logger"
	"myApi/interface/customer"
	"myApi/interface/product"
	productcharacteristics "myApi/interface/product/productCharacteristics"
	saleitem "myApi/interface/saleItem"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/commissionByProduct.sql
var reportSQL string

type SaleContract struct {
	Id         int                       `json:"id"`
	CustomerId int                       `json:"customer_id"`
	Customer   string                    `json:"customer"`
	SaleValue  float64                   `json:"sale_value"`
	Status     string                    `json:"status"`
	Products   saleitem.SaleItemContract `json:"products"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
}

type SaleItensContract struct {
	Name                string  `json:"name"`
	SaleValue           float64 `json:"sale_value"`
	Commission          float64 `json:"commission_by_produtc"`
	CommissionGenerated float64 `json:"commission_generated"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (s SaleContract) Validate() map[string]string {
	errorsField := make(map[string]string)
	var subTotal float64

	for _, p := range s.Products {
		subTotal += calchelper.CalculateTotalSale(p.SaleValue, p.Qtde)
	}

	if subTotal <= 0 {
		errorsField["sub_total"] = "O valor da venda não pode ser zerado."
	}

	if _, err := customer.Show(s.CustomerId); err != nil {
		u.ErrorLogger.Println("Erro no select do cliente para validar a venda: ", err)
		errorsField["customer_id"] = fmt.Sprintf("%s", err)
	}

	if err := s.Products.Validate(); len(err) > 0 {
		errorsField["products"] = fmt.Sprintf("%s", err)
	}

	u.GeneralLogger.Println("SubTotal da venda aqui: ", subTotal)

	return errorsField
}

func GetAll() ([]SaleContract, error) {
	var sales []SaleContract

	query := `
		SELECT
			id,
			customer_id,
			customer,
			sale_value,
			status
			
		FROM
			sales
	`

	rows, err := conn.Query(
		ctx,
		query,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s SaleContract

		if err := rows.Scan(
			&s.Id,
			&s.CustomerId,
			&s.Customer,
			&s.SaleValue,
			&s.Status,
		); err != nil {
			u.GeneralLogger.Println("Erro: ", err)
			return nil, err
		}

		sales = append(sales, s)
	}

	return sales, nil
}

func Show(id int) (*SaleContract, error) {
	var s SaleContract

	query := `
		SELECT
			id,
			customer_id,
			customer,
			sale_value,
			status
			
		FROM
			sales

		WHERE
			id = $1
	`

	if err := conn.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&s.Id,
		&s.CustomerId,
		&s.Customer,
		&s.SaleValue,
		&s.Status,
	); err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da venda")
		return nil, err
	}

	queryFromItens := `
		SELECT
			id,
			sale_id,
			product_id,
			name,
			qtde,
			sale_value,
			status
		FROM
			sale_itens

		WHERE
			sale_id = $1
	`

	rows, err := conn.Query(
		ctx,
		queryFromItens,
		id,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao fazer o select")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var productCharacteristic productcharacteristics.ProductCharacteristicsContract

		var p struct {
			Id                         int                                                    `json:"id"`
			SaleId                     int                                                    `json:"sale_id"`
			ProductId                  int                                                    `json:"product_id"`
			Name                       string                                                 `json:"name"`
			Qtde                       int                                                    `json:"qtde"`
			SaleValue                  float64                                                `json:"price"`
			Status                     string                                                 `json:"status"`
			ProductWithCharacteristics *productcharacteristics.ProductCharacteristicsContract `json:"product_with_characteristics"`
			CreatedAt                  time.Time                                              `json:"created_at"`
			UpdatedAt                  time.Time                                              `json:"cpdated_at"`
		}

		if err := rows.Scan(
			&p.Id,
			&p.SaleId,
			&p.ProductId,
			&p.Name,
			&p.Qtde,
			&p.SaleValue,
			&p.Status,
		); err != nil {
			u.ErrorLogger.Printf("Erro no select dos itens da venda - %s", err)
			return nil, err
		}

		err := conn.QueryRow(
			ctx,
			`
				SELECT
					product_id,
					sale_id,
					product_grid_id,
					size_saled,
					grid_qtde

				FROM
					sale_itens_grid

				WHERE
					sale_id = $1 AND
					product_id = $2
			`,
			id,
			p.ProductId,
		).Scan(
			&productCharacteristic.ProductId,
			&productCharacteristic.SaleId,
			&productCharacteristic.Id,
			&productCharacteristic.Size,
			&productCharacteristic.GridQtde,
		)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			u.ErrorLogger.Println("Erro ao pegar os dados da sale_itens_grid:", err)
			return nil, err
		}

		if errors.Is(err, pgx.ErrNoRows) {
			p.ProductWithCharacteristics = nil
		} else {
			p.ProductWithCharacteristics = &productCharacteristic
		}

		s.Products = append(s.Products, p)
	}

	return &s, nil
}

func ShowTotalCommission(id int) (*[]SaleItensContract, error) {
	var saleItens []SaleItensContract

	rows, err := conn.Query(
		ctx,
		string(reportSQL),
		id,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao executar a query:", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var sl SaleItensContract

		if err := rows.Scan(
			&sl.Name,
			&sl.SaleValue,
			&sl.Commission,
			&sl.CommissionGenerated,
		); err != nil {
			u.ErrorLogger.Println("Erro ao ler os dados da query:", err)
			return nil, err
		}

		saleItens = append(saleItens, sl)
	}

	return &saleItens, nil
}

func (s *SaleContract) Create() (int, error) {
	tx, err := conn.Begin(ctx)

	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)

	u.InfoLogger.Println("ID cliente: ", s.CustomerId)

	if s.CustomerId > 1 && s.Customer != "Consumidor padrão" {
		u.InfoLogger.Println("Cliente diferente do padrão")
		otherCustomer, err := customer.Show(s.CustomerId)

		if err != nil {
			u.ErrorLogger.Println("Erro ao localizar o cliente:", err)
			return 0, err
		}

		s.Customer = otherCustomer.Name
		s.CustomerId = otherCustomer.Id
	}

	var saleId int

	if err = tx.QueryRow(
		ctx,
		`
			INSERT INTO sales
				(customer_id, customer, sale_value)

			VALUES
				($1, $2, $3)

			RETURNING 
				id
		`,
		s.CustomerId,
		s.Customer,
		s.SaleValue,
	).Scan(
		&saleId,
	); err != nil {
		u.ErrorLogger.Println("Erro no create (sale-contract): ", err)
		return 0, err
	}

	s.Id = saleId

	for _, itemSale := range s.Products {
		productData, err := product.Show(itemSale.ProductId)

		if err != nil {
			u.ErrorLogger.Printf("Erro ao localizar o produto pelo ID %d, %s", itemSale.ProductId, err)
			return 0, err
		}

		totalSale := calchelper.CalculateTotalSale(itemSale.SaleValue, itemSale.Qtde)

		if !productData.UseGrid && itemSale.ProductWithCharacteristics == nil {
			u.InfoLogger.Println("O produto da venda não tem grade, vai gravar normal")

			if _, err := tx.Exec(
				ctx,
				`
					INSERT INTO sale_itens
						(product_id, name, qtde, sale_value, sale_id)

					VALUES
						($1, $2, $3, $4, $5)
						
					RETURNING 
						id,
						name, 
						status
				`,
				itemSale.ProductId,
				productData.Name,
				itemSale.Qtde,
				totalSale,
				saleId,
			); err != nil {
				u.ErrorLogger.Println("Erro no insert no sale_itens:", err)
				return 0, err
			}

			if err = productData.DiscountedQtde(ctx, tx, itemSale.Qtde); err != nil {
				u.ErrorLogger.Println("Erro no insert no DiscountedQtde:", err)
				return 0, err
			}
		}

		if productData.UseGrid && itemSale.ProductWithCharacteristics != nil {
			u.InfoLogger.Println("O produto da venda tem grade, vai gravar na sale_itens_grid também")

			if err = tx.QueryRow(
				ctx,
				`
					INSERT INTO sale_itens
						(product_id, name, qtde, sale_value, sale_id)

					VALUES
						($1, $2, $3, $4, $5)
						
					RETURNING 
						id,
						name,
						status
				`,
				itemSale.ProductId,
				productData.Name,
				itemSale.Qtde,
				totalSale,
				saleId,
			).Scan(
				&itemSale.Id,
				&itemSale.Name,
				&itemSale.Status,
			); err != nil {
				u.ErrorLogger.Println("Erro no insert no sale_itens:", err)
				return 0, err
			}

			itemSale.SaleId = saleId

			u.InfoLogger.Printf("Vai descontar a qtde %d da grade %s do produto", itemSale.Qtde, itemSale.ProductWithCharacteristics.Size)
			if err = itemSale.ProductWithCharacteristics.DiscountedGridQtde(ctx, tx, itemSale.Qtde, productcharacteristics.Size(itemSale.ProductWithCharacteristics.Size)); err != nil {
				u.ErrorLogger.Println("Erro ao descontar a qtde da grade do produto:", err)
				return 0, err
			}

			u.InfoLogger.Println("Vai gravar na sale_itens_grid também:", itemSale.ProductWithCharacteristics)

			if itemSale.ProductWithCharacteristics == nil {
				u.ErrorLogger.Println("O produto possui grade porém está com as caracteristicas ausentes:", itemSale)
				return 0, fmt.Errorf("O produto possui grade porém está com as caracteristicas ausentes.")
			}

			if _, err = tx.Exec(
				ctx,
				`
					INSERT INTO sale_itens_grid
						(product_id, sale_id, product_grid_id, size_saled, grid_qtde)
	
					VALUES
						($1, $2, $3, $4, $5)
					
				`,
				itemSale.ProductId,
				saleId,
				itemSale.ProductWithCharacteristics.Id,
				itemSale.ProductWithCharacteristics.Size,
				itemSale.Qtde,
			); err != nil {
				u.ErrorLogger.Println("Erro no insert no sale_itens_grid:", err)
				return 0, err
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return saleId, err
}

func (s *SaleContract) InsertNewItens() error {
	if s.Id <= 0 {
		return fmt.Errorf("Identificador da venda inválido.")
	}

	u.InfoLogger.Println("InsertNewItens started")

	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transição: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	for _, p := range s.Products {
		u.InfoLogger.Println("Produto: ", p)
		u.InfoLogger.Println("ProductWithCharacteristics: ", p.ProductWithCharacteristics)

		var useGrid bool
		var oldQtdeSaleItem int

		if err := tx.QueryRow(
			ctx,
			`
				SELECT
					use_grid
				FROM
					products
				WHERE
					id = $1
			`,
			p.ProductId,
		).Scan(
			&useGrid,
		); err != nil {
			u.ErrorLogger.Println("Erro ao conferir se o produto usa grade:", err)
			return err
		}

		if useGrid {
			if p.ProductWithCharacteristics == nil {
				u.ErrorLogger.Println("Produto utiliza grade mas a grade está ausente.", p)
				return fmt.Errorf("Produto utiliza grade mas a grade está ausente.")
			}

			u.InfoLogger.Println("O produto usa grade")
			u.InfoLogger.Println("1 - Vai atualizar a qtde do item na venda")

			if _, err := tx.Exec(
				ctx,
				`
					UPDATE
						sale_itens
	
					SET
						qtde = $1
	
					WHERE
						sale_id = $2 AND
						product_id = $3
				`,
				p.Qtde,
				s.Id,
				p.ProductId,
			); err != nil {
				u.ErrorLogger.Println("Erro ao atualizar a qtde do item: ", err)
				return err
			}

			u.InfoLogger.Println("2 - Vai atualizar a qtde da grade do item na venda")

			if _, err := tx.Exec(
				ctx,
				`
					UPDATE
						sale_itens_grid
	
					SET
						grid_qtde = $1
	
					WHERE
						sale_id = $2 AND
						product_id = $3 AND
						size_saled = $4
				`,
				p.Qtde,
				s.Id,
				p.ProductId,
				p.ProductWithCharacteristics.Size,
			); err != nil {
				u.ErrorLogger.Println("Erro ao atualizar a grade do item: ", err)
				return err
			}

			u.InfoLogger.Println("Vai pegar a qtde que foi passada na venda")

			err := tx.QueryRow(
				ctx,
				`
					SELECT
						grid_qtde
	
					FROM
						sale_itens_grid
	
					WHERE
						size_saled = $1 AND
						product_id = $2 AND
						sale_id = $3
	
				`,
				p.ProductWithCharacteristics.Size,
				p.ProductId,
				s.Id,
			).Scan(
				&oldQtdeSaleItem,
			)

			if errors.Is(err, pgx.ErrNoRows) {
				u.InfoLogger.Println("Não possui essa grade inserida, vai inserir uma nova grade")
				if _, err = tx.Exec(
					ctx,
					`
						INSERT INTO sale_itens_grid
							(product_id, sale_id, product_grid_id, size_saled, grid_qtde)

						VALUES
							($1, $2, $3, $4, $5)
						
					`,
					p.ProductId,
					s.Id,
					p.ProductWithCharacteristics.Id,
					p.ProductWithCharacteristics.Size,
					p.Qtde,
				); err != nil {
					u.ErrorLogger.Println("Erro no insert no sale_itens_grid:", err)
					return err
				}
			}

			if err != nil {
				u.ErrorLogger.Println("Erro ao localizar a qtde anterior da grade do produto da venda:", err)
				return err
			}

			u.InfoLogger.Println("Qtde que foi passada na venda: ", oldQtdeSaleItem)

			if err == pgx.ErrNoRows {
				u.InfoLogger.Println("Não foi localizado qtde anterior")
				oldQtdeSaleItem = 0
			}

			u.InfoLogger.Println("p.ProductWithCharacteristics.GridQtde:", p.ProductWithCharacteristics.GridQtde)
			diff := p.ProductWithCharacteristics.GridQtde - oldQtdeSaleItem

			u.InfoLogger.Println("Diferença:", diff)
			u.InfoLogger.Println("Vai alterar a qtde passada na venda.")

			if diff != 0 {
				u.InfoLogger.Printf("Vai alterar a qtde da grade do produto: id - %d, diferença: %d, grade - %s", p.ProductId, diff, p.ProductWithCharacteristics.Size)
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							product_grids	
						SET 
							grid_qtde = grid_qtde - $2
						WHERE
							product_id = $1 AND
							size = $3
					`,
					p.ProductId,
					diff,
					p.ProductWithCharacteristics.Size,
				); err != nil {
					u.ErrorLogger.Println("Erro ao fazer o update na grade do produto: ", err)
					return err
				}

				u.InfoLogger.Printf("Vai alterar a qtde do produto: id - %d", p.ProductId)
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							products
						SET 
							qtde = (
								SELECT
									COALESCE(SUM(grid_qtde), 0)

								FROM
									product_grids

								WHERE
									product_id = $1
							)
						WHERE
							id = $1
					`,
					p.ProductId,
				); err != nil {
					u.ErrorLogger.Println("Erro ao fazer o update na grade do produto: ", err)
					return err
				}
			}
		} else {
			if err := tx.QueryRow(
				ctx,
				`
					SELECT 
						qtde
	
					FROM
						sale_itens
	
					WHERE
						sale_id = $1 AND
						product_id = $2
	
				`,
				s.Id,
				p.ProductId,
			).Scan(&oldQtdeSaleItem); err != nil && err != pgx.ErrNoRows {
				u.ErrorLogger.Println("Erro ao localizar a qtde anterior do produto da venda")
				return err
			}

			if err == pgx.ErrNoRows {
				oldQtdeSaleItem = 0
			}

			diff := p.Qtde - oldQtdeSaleItem

			var itenSaleExisit bool

			// CASE FOR CONFIRM IF PRODUCTS EXISTI IN SALE
			if err := tx.QueryRow(
				ctx,
				`	
					SELECT
						id

					FROM
						sale_itens

					WHERE
						product_id = $1 AND 
						sale_id = $2

					LIMIT
						1
				`,
				p.ProductId,
				s.Id,
			).Scan(&itenSaleExisit); err != nil && !errors.Is(err, pgx.ErrNoRows) {
				u.ErrorLogger.Println("Ocorroe um erro ao salvar a venda:", err)
				return err
			}
			// --------------------------------------------------------- \\

			if !itenSaleExisit {
				// CASE FOR INSERT NEW ITENS ONLY DON'T LOCALIZATE ID PRODUCT
				if _, err := tx.Exec(
					ctx,
					`
						INSERT INTO sale_itens
							(product_id, name, qtde, sale_value, sale_id)
		
						VALUES
							($1, $2, $3, $4, $5)
					`,
					p.ProductId,
					p.Name,
					p.Qtde,
					p.SaleValue,
					s.Id,
				); err != nil {
					u.ErrorLogger.Println("Erro ao inserir/atualizar item: ", err)
					return err
				}
				// --------------------------------------------------------- \\
			} else {
				// CASE FOR UPDATE ITENS
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							sale_itens

						SET
							qtde = $1,
							sale_value = $2

						WHERE
							product_id = $3 AND
							sale_id = $4
					`,
					p.Qtde,
					p.SaleValue,
					p.ProductId,
					s.Id,
				); err != nil {
					u.ErrorLogger.Println("Erro ao inserir/atualizar item: ", err)
					return err
				}
				// --------------------------------------------------------- \\
			}

			if diff != 0 {
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							products	
						SET 
							qtde = qtde - $2
						WHERE
							id = $1
					`,
					p.ProductId,
					diff,
				); err != nil {
					u.ErrorLogger.Println("Erro ao fazer o update no estoque: ", err)
					return err
				}
			}
		}
	}

	u.InfoLogger.Println("Vai alterar o valor total da venda")
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE 
				sales
			SET
				sale_value = (
					SELECT
						COALESCE(SUM(qtde * sale_value), 0)

					FROM
						sale_itens

					WHERE
						sale_id = $1
				)
			WHERE
				id = $1
		`,
		s.Id,
	); err != nil {
		u.ErrorLogger.Println("Erro ao alterar o total da venda depois da inserção/alteração dos itens: ", err)
		return err
	}

	u.InfoLogger.Println("Valor da venda atualizado.")

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao comitar: ", err)

		return err
	}

	return nil
}
