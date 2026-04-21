package processpayment

import (
	"context"
	"errors"
	"fmt"
	u "myApi/helpers/logger"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	"myApi/interface/product"
	"myApi/interface/sale"
	"myApi/interface/shopping"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

type PayMentBody struct {
	SpecieId   int     `json:"id"`
	Specie     string  `json:"specie"`
	AmountPaid float64 `json:"amount"`
}

type PayContract struct {
	SaleId     int           `json:"sale_id"`
	ShoppingId int           `json:"shopping_id"`
	Species    []PayMentBody `json:"species"`
}

type CancelContract struct {
	SaleId     int `json:"sale_id"`
	ShoppingId int `json:"shopping_id"`
}

func (payMent PayContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	if len(payMent.Species) == 0 {
		errorsField["species"] = "Pagamento ausente."
		return errorsField
	}

	if payMent.ShoppingId > 0 && payMent.SaleId > 0 {
		u.ErrorLogger.Printf("Identificadores de compra ou venda não podem ser definidos ao mesmo tempo, shopping_id: %d, sale_id: %d", payMent.ShoppingId, payMent.SaleId)

		errorsField["id"] = "Identificadores de compra ou venda não podem ser definidos ao mesmo tempo."
	}

	if payMent.ShoppingId == 0 && payMent.SaleId > 0 {
		saleData, err := sale.Show(payMent.SaleId)

		if err != nil {
			errorsField["database"] = fmt.Sprintf("Ocorreu um erro ao conferir se a venda existe: %s", err)
			return errorsField
		}

		if saleData == nil {
			errorsField["sale_id"] = fmt.Sprintf("O identificador da venda está incorreto, %s", err)
			return errorsField
		}
	}

	var totalPaide float64

	if payMent.SaleId == 0 && payMent.ShoppingId > 0 {
		shoppingData, err := shopping.Show(payMent.ShoppingId)

		if err != nil {
			errorsField["database"] = fmt.Sprintf("Ocorreu um erro ao conferir se a compra existe: %s", err)
			return errorsField
		}

		if shoppingData == nil {
			errorsField["shopping_id"] = fmt.Sprintf("O identificador da compra está incorreto, %s", err)
			return errorsField
		}
	}

	for _, pay := range payMent.Species {
		if pay.Specie != "Dinheiro" && pay.Specie != "Pix" {
			errorsField["species.specie"] = "A espécie de pagamento precisa ser Dinheiro ou Pix."
			return errorsField
		}

		u.GeneralLogger.Println("Forma de pagamento aqui: ", payMent)

		totalPaide += pay.AmountPaid
	}

	if totalPaide <= 0 {
		errorsField["amount_paid"] = "O pagamento não pode ser menor que zero."
		return errorsField
	}

	return errorsField
}

func createInCashRegister(tx pgx.Tx, inputValue, outputValue float64, saleId, shoppingId int, customer customer.CustomerContract, specie PayMentBody,
) map[string]string {
	errorsField := make(map[string]string)
	var c cashRegister.CashRegisterContract

	c.SpecieId = specie.SpecieId
	c.Specie = specie.Specie

	c.SaleId = saleId
	c.ShoppingId = shoppingId

	c.CustomerId = customer.Id
	c.Customer = customer.Name

	if inputValue > 0 && outputValue > 0 {
		u.ErrorLogger.Println("Um registro no caixa não pode ter um valor de entrada e um de saída no mesmo registro.")

		errorsField["input_value"] = "Um registro no caixa não pode ter um valor de entrada no mesmo registro de uma saída."
		errorsField["output_value"] = "Um registro no caixa não pode ter um valor de saída no mesmo registro de uma entrada."

		return errorsField
	}

	if inputValue > 0 {
		if err := c.Create(tx, inputValue, 0.0); len(err) > 0 {
			return err
		}
	}

	if outputValue > 0 {
		if err := c.Create(tx, 0.0, outputValue); len(err) > 0 {
			return err
		}
	}

	return errorsField
}

func PayMentShoppingOrSale(payMent PayContract, totalPaide float64) error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transiction no PayMentShoppingOrSale: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	if payMent.SaleId > 0 {
		if err := processSalePayment(ctx, tx, payMent); err != nil {
			u.ErrorLogger.Println("Erro ao processar o pagamento da venda:", err)
			return err
		}
	} // Processo para pagamento de VENDA

	if payMent.ShoppingId > 0 {
		if err := processShoppingPayment(ctx, tx, payMent); err != nil {
			u.ErrorLogger.Println("Erro ao processar o pagamento da compra:", err)
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro no commit do PayMentShoppingOrSale: ", err)
		return err
	}

	return nil
}

func getShoppingData(ctx context.Context, tx pgx.Tx, shoppingID int) (shopping.ShoppingContract, error) {
	var s shopping.ShoppingContract

	if err := tx.QueryRow(
		ctx,
		`
			SELECT
				total_shopping,
				status

			FROM
				shopping
			
			WHERE
				id = $1
		`,
		shoppingID,
	).Scan(
		&s.TotalShopping,
		&s.Status,
	); err != nil {
		u.ErrorLogger.Println("Erro ao localizar a compra: ", err)
		return shopping.ShoppingContract{}, err
	}

	return s, nil
}

func getSaleData(ctx context.Context, tx pgx.Tx, saleID int) (sale.SaleContract, error) {
	var s sale.SaleContract

	if err := tx.QueryRow(
		ctx,
		`
			SELECT
				customer_id,
				sale_value,
				status
			FROM
				sales
			WHERE
				id = $1
		`,
		saleID,
	).Scan(
		&s.CustomerId,
		&s.SaleValue,
		&s.Status,
	); err != nil {
		u.ErrorLogger.Println("Erro no select da venda:", err)
		return sale.SaleContract{}, err
	}

	return s, nil
}

func processSalePayment(ctx context.Context, tx pgx.Tx, payMent PayContract) error {
	s, err := getSaleData(ctx, tx, payMent.SaleId)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da venda:", err)
		return err
	}

	if s.Status == "Concluída" {
		return fmt.Errorf("Essa venda já está finalizada.")
	}

	u.InfoLogger.Println("Vai fazer o insert no sale_pay_ment")

	if err := insertSalePayMents(ctx, tx, payMent.SaleId, s.CustomerId, payMent); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o inset no sale_pay_ment:", err)
		return err
	}

	u.SuccessLoger.Println("Pagamento inserido, vai processar o estoque.")

	if err := processSaleStock(ctx, tx, payMent.SaleId, true); err != nil {
		u.ErrorLogger.Println("Erro ao processar o estoque:", err)
		return err
	}

	u.SuccessLoger.Println("Estoque processado, vai finalizar a venda e os itens.")

	if err := finallySale(ctx, tx, payMent.SaleId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar a venda:", err)
		return err
	}

	u.SuccessLoger.Println("Venda finalizada, vai finalizar os itens da venda.")

	if err := finallySaleItens(ctx, tx, payMent.SaleId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar os itens da venda:", err)
		return err
	}

	u.SuccessLoger.Println("Itens da venda finalizada, vai finalizar a grade dos itens da venda.")

	if err := finallySaleGridItens(ctx, tx, payMent.SaleId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar a grade dos itens da venda:", err)
		return err
	}

	u.SuccessLoger.Println("Venda finalizada.")

	return nil
}

func finallySale(ctx context.Context, tx pgx.Tx, saleID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE 
				sales
			SET
				status = 'Concluída'		
			WHERE 
				id = $1
		`,
		saleID,
	); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar a venda: ", err)
		return err
	}

	return nil
}

func finallySaleItens(ctx context.Context, tx pgx.Tx, saleID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				sale_itens
			SET
				status = 'Concluída'	
			WHERE 
				sale_id = $1
		`,
		saleID,
	); err != nil {
		u.ErrorLogger.Println("Erro no update dos itens da venda para Concluída: ", err)
		return err
	}
	return nil
}

func finallySaleGridItens(ctx context.Context, tx pgx.Tx, saleID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				sale_itens
			SET
				status = 'Concluída'	
			WHERE 
				sale_id = $1
		`,
		saleID,
	); err != nil {
		u.ErrorLogger.Println("Erro no update dos itens da venda para Concluída: ", err)
		return err
	}

	return nil
}

func insertSalePayMents(ctx context.Context, tx pgx.Tx, saleID, customerID int, species PayContract) error {
	for _, specie := range species.Species {
		if specie.AmountPaid > 0 {
			if _, err := tx.Exec(
				ctx,
				`
					INSERT INTO sale_pay_ment
						(sale_id, specie_id, specie, amount_paid)
		
					VALUES
						($1, $2, $3, $4)
				`,
				saleID,
				specie.SpecieId,
				specie.Specie,
				specie.AmountPaid,
			); err != nil {
				u.ErrorLogger.Println("Erro no insert do sale_pay_ment no paySale: ", err)
				return err

			}

			c, err := customer.Show(customerID)

			if err != nil {
				u.ErrorLogger.Println("Erro ao pegar os dados do cliente: ", err)
				return err
			}

			if err := createInCashRegister(
				tx,
				specie.AmountPaid,
				0.0,
				saleID,
				0,
				*c,
				specie,
			); len(err) > 0 {
				return fmt.Errorf("Erros: %s", err)
			}
		}
	}

	return nil
}

func processSaleStock(ctx context.Context, tx pgx.Tx, saleID int, isToRemove bool) error {
	var saleItensGridId int

	if err := tx.QueryRow(
		ctx,
		`
			SELECT
				id
			FROM
				sale_itens_grid				
			WHERE
				sale_id = $1
			LIMIT 1
		`,
		saleID,
	).Scan(&saleItensGridId); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao conferir se existe itens com grade na compra:", err)
		return err
	}

	if saleItensGridId < 0 {
		u.InfoLogger.Println("A venda não possui grade")

		saleItensRows, err := tx.Query(
			ctx,
			`
				SELECT
					product_id,
					qtde
				FROM
					sale_itens
				WHERE
					status = 'Concluída' AND
					sale_id = $1
			`,
			saleID,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao pegar os itens da venda:", err)
			return err
		}

		defer saleItensRows.Close()

		for saleItensRows.Next() {
			var saleItem sale.SaleItensContract

			if err := saleItensRows.Scan(
				&saleItem.ProductId,
				&saleItem.Qtde,
			); err != nil {
				u.ErrorLogger.Println("Erro ao ler os itens da venda:", err)
				return err
			}

			product := product.ProductContract{
				Id: saleItem.ProductId,
			}

			if isToRemove {
				if err := product.DiscountedQtde(ctx, tx, saleItem.Qtde, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao descontar a qtde do item ao estoque:", err)
					return err
				}
			} else {
				if err := product.AddQtde(ctx, tx, saleItem.Qtde, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao adicionar a qtde do item ao estoque:", err)
					return err
				}
			}
		}

	} else {
		u.InfoLogger.Println("A venda possui grade")

		saleItensGridRows, err := tx.Query(
			ctx,
			`
				SELECT
					product_id,
					grid_qtde,
					size_saled

				FROM
					sale_itens_grid
				WHERE
					sale_id = $1 AND 
					status = 'Concluída'
			`,
			saleID,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao pegar a grade dos itens da venda:", err)
			return err
		}

		defer saleItensGridRows.Close()

		for saleItensGridRows.Next() {
			var saleItemGrids product.ProductGrids

			if err := saleItensGridRows.Scan(
				&saleItemGrids.ProductId,
				&saleItemGrids.GridQtde,
				&saleItemGrids.Size,
			); err != nil {
				u.ErrorLogger.Println("Erro ao pegar a grade dos itens da compra:", err)
				return err
			}

			productDataId := product.ProductContract{
				Id: saleItemGrids.ProductId,
			}

			productGrid := product.ProductGrids{
				ProductId: productDataId.Id,
				Size:      saleItemGrids.Size,
				GridQtde:  saleItemGrids.GridQtde,
			}

			if isToRemove {
				if err := productDataId.DiscountedQtde(ctx, tx, 0, true, &productGrid); err != nil {
					u.ErrorLogger.Println("Erro ao descontar a qtde da grade do item ao estoque:", err)
					return err
				}
			} else {
				if err := productDataId.AddQtde(ctx, tx, 0, true, &productGrid); err != nil {
					u.ErrorLogger.Println("Erro ao adicionar a qtde da grade do item ao estoque:", err)
					return err
				}
			}
		}
	}

	return nil
}

func processShoppingPayment(ctx context.Context, tx pgx.Tx, payMent PayContract) error {
	var s shopping.ShoppingContract

	if err := tx.QueryRow(
		ctx,
		`
			SELECT
				id,
				load,
				operation,
				status
			FROM			
				shopping

			WHERE
				id = $1
		`,
		payMent.ShoppingId,
	).Scan(
		&s.Id,
		&s.Load,
		&s.Operation,
		&s.Status,
	); err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da compra:", err)
		return err
	}

	if s.Status == "Concluída" {
		return fmt.Errorf("Essa compra já está finalizada.")
	}

	u.InfoLogger.Println("Vai fazer o insert no sale_pay_ment")

	if err := insertShoppingPayMents(ctx, tx, payMent.ShoppingId, payMent); err != nil {
		return err
	}

	u.InfoLogger.Println("Pagamento inserido, vai processar o estoque.")

	if err := processShoppingStock(ctx, tx, payMent.ShoppingId, true); err != nil {
		u.ErrorLogger.Println("Erro ao processar o estoque:", err)
		return err
	}

	u.InfoLogger.Println("Todos os processos foram realizados, vai finalizar a compra e os itens.")

	if err := finallyShopping(ctx, tx, payMent.ShoppingId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar a compra:", err)
		return err
	}

	if err := finallyShoppingItens(ctx, tx, payMent.ShoppingId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar os itens da compra:", err)
		return err
	}

	if err := finallyShoppingGridItens(ctx, tx, payMent.ShoppingId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar a grade dos itens da compra:", err)
		return err
	}

	return nil
}

func processShoppingStock(ctx context.Context, tx pgx.Tx, shoppingID int, isToAdd bool) error {
	hasGrid := true
	var shoppingItensGridId int

	err := tx.QueryRow(
		ctx,
		`
			SELECT
				id
			FROM
				shopping_itens_grid
			WHERE
				shopping_id = $1
			LIMIT 1
		`,
		shoppingID,
	).Scan(&shoppingItensGridId)

	if errors.Is(err, pgx.ErrNoRows) {
		hasGrid = false

	} else if err != nil {
		u.ErrorLogger.Println("Erro ao consultar se o produto possui grade:", err)
		return err

	}

	if !hasGrid {
		u.InfoLogger.Println("A compra não possui grade")

		shoppingItensRows, err := tx.Query(
			ctx,
			`
				SELECT
					product_id,
					qtde_purchased,
					purchased_value
				FROM
					shopping_itens
					
				WHERE
					status = 'Concluída' AND
					shopping_id = $1
			`,
			shoppingID,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao pegar os itens da compra:", err)
			return err
		}

		defer shoppingItensRows.Close()

		for shoppingItensRows.Next() {
			var shoppingItem shopping.ShoppingItenContract

			if err := shoppingItensRows.Scan(
				&shoppingItem.ProductId,
				&shoppingItem.QtdePurchased,
			); err != nil {
				u.ErrorLogger.Println("Erro ao ler os itens da venda:", err)
				return err
			}

			product := product.ProductContract{
				Id: shoppingItem.ProductId,
			}

			if isToAdd {
				if err := product.AddQtde(ctx, tx, shoppingItem.QtdePurchased, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao adicionar a qtde do item ao estoque:", err)
					return err
				}

			} else {
				if err := product.DiscountedQtde(ctx, tx, shoppingItem.QtdePurchased, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao remover a qtde do item ao estoque:", err)
					return err
				}

			}
		}

	} else {
		u.InfoLogger.Println("A compra possui grade")

		shoppingItensGridRows, err := tx.Query(
			ctx,
			`
				SELECT
					product_id,
					size_saled,
					grid_qtde

				FROM
					shopping_itens_grid

				WHERE
					shopping_id = $1
			`,
			shoppingID,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao pegar a grade dos itens da venda:", err)
			return err
		}

		defer shoppingItensGridRows.Close()

		for shoppingItensGridRows.Next() {
			var shoppingItemGrids product.ProductGrids

			if err := shoppingItensGridRows.Scan(
				&shoppingItemGrids.ProductId,
				&shoppingItemGrids.GridQtde,
				&shoppingItemGrids.Size,
			); err != nil {
				u.ErrorLogger.Println("Erro ao pegar a grade dos itens da compra:", err)
				return err
			}

			productDataId := product.ProductContract{
				Id: shoppingItemGrids.ProductId,
			}

			productGrid := product.ProductGrids{
				ProductId: productDataId.Id,
				Size:      shoppingItemGrids.Size,
				GridQtde:  shoppingItemGrids.GridQtde,
			}

			if isToAdd {
				if err := productDataId.AddQtde(ctx, tx, 0, true, &productGrid); err != nil {
					u.ErrorLogger.Println("Erro ao adicionar a qtde da grade do item ao estoque:", err)
					return err
				}
			} else {
				if err := productDataId.DiscountedQtde(ctx, tx, 0, true, &productGrid); err != nil {
					u.ErrorLogger.Println("Erro ao remover a qtde da grade do item ao estoque:", err)
					return err
				}
			}
		}
	}

	return nil
}

func finallyShopping(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE 
				shopping
			SET
				status = 'Concluída'		
			WHERE 
				id = $1
		`,
		shoppingID,
	); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar a compra: ", err)
		return err
	}

	return nil
}

func finallyShoppingItens(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				shopping_itens
			SET
				status = 'Concluída'	
			WHERE 
				shopping_id = $1
		`,
		shoppingID,
	); err != nil {
		u.ErrorLogger.Println("Erro no update dos itens da compra para Concluída: ", err)
		return err
	}

	return nil
}

func finallyShoppingGridItens(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				shopping_itens_grid
			SET
				status = 'Concluída'	
			WHERE 
				shopping_id = $1
		`,
		shoppingID,
	); err != nil {
		u.ErrorLogger.Println("Erro no update da grade dos itens da compra para Concluída: ", err)
		return err
	}

	return nil
}

func insertShoppingPayMents(ctx context.Context, tx pgx.Tx, shoppingID int, species PayContract) error {
	for _, specie := range species.Species {
		if specie.AmountPaid > 0 {
			if _, err := tx.Exec(
				ctx,
				`
					INSERT INTO shopping_pay_ment
						(shopping_id, specie_id, specie, amount_paid)
		
					VALUES
						($1, $2, $3, $4)
					
					RETURNING
						id
				`,
				shoppingID,
				specie.SpecieId,
				specie.Specie,
				specie.AmountPaid,
			); err != nil {
				u.ErrorLogger.Println("Erro no insert do shopping_pay_ment: ", err)
				return err

			}

			c, err := customer.Show(1)

			if err != nil {
				u.ErrorLogger.Println("Erro ao pegar os dados do cliente: ", err)
				return err
			}

			if err := createInCashRegister(
				tx,
				specie.AmountPaid,
				0.0,
				0,
				shoppingID,
				*c,
				specie,
			); len(err) > 0 {
				return fmt.Errorf("Erros: %s", err)
			}
		}
	}

	return nil
}

func processSaleCancel(ctx context.Context, tx pgx.Tx, saleID int) error {
	u.InfoLogger.Println("Called processSaleCancel")
	s, err := getSaleData(ctx, tx, saleID)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da venda:", err)
		return err
	}

	switch s.Status {
	case "Cancelada":
		u.ErrorLogger.Printf("Essa venda n° %d já está cancelada", s.Id)
		return fmt.Errorf("Essa venda n° %d já está cancelada", s.Id)

	case "Pendente":
		if err := cancelSale(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a venda:", err)
			return err
		}

		u.SuccessLoger.Println("Venda cancelada, vai cancelar os itens")

		if err := cancelSaleItens(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar os itens da venda:", err)
			return err
		}

		u.SuccessLoger.Println("Itens da venda cancelados, vai cancelar a grade dos itens")

		if err := cancelSaleGridItens(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a grade dos itens da venda:", err)
			return err
		}

		u.SuccessLoger.Println("Venda cancelada.")

	case "Concluída":
		if err := backSalePayMent(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao realizar o estorno do caixa:", err)
			return err
		}

		u.SuccessLoger.Println("Dinheiro estornado, vai cancelar a venda.")

		if err := cancelSale(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a venda:", err)
			return err
		}

		if err := cancelSaleItens(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar os itens da venda:", err)
			return err
		}

		u.SuccessLoger.Println("Itens da venda cancelados, vai cancelar a grade dos itens")

		if err := cancelSaleGridItens(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a grade dos itens da venda:", err)
			return err
		}

		u.SuccessLoger.Println("Venda cancelada.")
	}

	return nil
}

func cancelSale(ctx context.Context, tx pgx.Tx, saleID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				sales
			SET
				status = 'Cancelada'
			WHERE
				id = $1
		`,
		saleID,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update sale para cancelado - %s", err)
		return err
	}

	return nil
}

func cancelSaleItens(ctx context.Context, tx pgx.Tx, saleID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				sale_itens

			SET
				status = 'Cancelada'

			WHERE
				sale_id = $1
		`,
		saleID,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
		return err
	}

	return nil
}

func cancelSaleGridItens(ctx context.Context, tx pgx.Tx, saleID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				sale_itens_grid
			SET
				status = 'Cancelada'
			WHERE
				sale_id = $1
		`,
		saleID,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
		return err
	}

	return nil
}

func backSalePayMent(ctx context.Context, tx pgx.Tx, saleID int) error {
	s, err := getSaleData(ctx, tx, saleID)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da venda:", err)
		return err
	}

	payMentFormsSelect, err := tx.Query(
		ctx,
		`
			SELECT
				specie_id,
				specie,
				amount_paid
			FROM
				sale_pay_ment
	
			WHERE
				sale_id = $1
		`,
		saleID,
	)

	if err != nil {
		u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
		return err
	}

	customer, err := customer.Show(s.CustomerId)

	if err != nil {
		u.ErrorLogger.Println("Erro no select do cliente para validar a venda: ", err)
		return err
	}

	defer payMentFormsSelect.Close()

	for payMentFormsSelect.Next() {
		var pf PayMentBody

		if err := payMentFormsSelect.Scan(
			&pf.SpecieId,
			&pf.Specie,
			&pf.AmountPaid,
		); err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
			return err
		}

		if err := createInCashRegister(
			tx,
			0.0,
			pf.AmountPaid,
			saleID,
			0,
			*customer,
			pf,
		); len(err) > 0 {
			u.ErrorLogger.Printf("Erro no insert de estorno no caixa venda - %s", err)
			return fmt.Errorf("Erro ao registrar o estorno no caixa.")
		}
	}

	return nil
}

func processShoppingCancel(ctx context.Context, tx pgx.Tx, saleID int) error {

	return nil
}

func cancelShopping(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				shopping
			SET
				status = 'Cancelada'
			WHERE
				id = $1
		`,
		shoppingID,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update sale para cancelado - %s", err)
		return err
	}

	return nil
}
func cancelShoppingGrid(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				shopping_itens

			SET
				status = 'Cancelada'

			WHERE
				sale_id = $1
		`,
		shoppingID,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
		return err
	}
	return nil
}

func cancelShoppingGridItens(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				sale_itens_grid
			SET
				status = 'Cancelada'
			WHERE
				sale_id = $1
		`,
		shoppingID,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
		return err
	}

	return nil
}

func backShoppingPayMent(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	s, err := getSaleData(ctx, tx, shoppingID)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da venda:", err)
		return err
	}

	payMentFormsSelect, err := tx.Query(
		ctx,
		`
			SELECT
				specie_id,
				specie,
				amount_paid
			FROM
				shopping_pay_ment
	
			WHERE
				sale_id = $1
		`,
		shoppingID,
	)

	if err != nil {
		u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
		return err
	}

	customer, err := customer.Show(s.CustomerId)

	if err != nil {
		u.ErrorLogger.Println("Erro no select do cliente para validar a venda: ", err)
		return err
	}

	defer payMentFormsSelect.Close()

	for payMentFormsSelect.Next() {
		var pf PayMentBody

		if err := payMentFormsSelect.Scan(
			&pf.SpecieId,
			&pf.Specie,
			&pf.AmountPaid,
		); err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
			return err
		}

		if err := createInCashRegister(
			tx,
			0.0,
			pf.AmountPaid,
			0,
			shoppingID,
			*customer,
			pf,
		); len(err) > 0 {
			u.ErrorLogger.Printf("Erro no insert de estorno no caixa venda - %s", err)
			return fmt.Errorf("Erro ao registrar o estorno no caixa.")
		}
	}

	return nil
}

func CancelSaleOrShopping(c CancelContract) error {
	u.InfoLogger.Println("called CancelSaleOrShopping")

	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Printf("Erro ao iniciar a transação - %s", err)
		return err
	}

	defer tx.Rollback(ctx)

	if c.SaleId > 0 {
		u.InfoLogger.Println("Cancelando uma venda")

	}

	if c.ShoppingId > 0 {
		u.InfoLogger.Println("Cancelando uma compra")

		var s shopping.ShoppingContract

		if err := tx.QueryRow(
			ctx,
			`
				SELECT
					total_shopping,
					status

				FROM
					shopping
				
				WHERE
					id = $1
			`,
			c.ShoppingId,
		).Scan(
			&s.TotalShopping,
			&s.Status,
		); err != nil {
			u.ErrorLogger.Println("Erro ao localizar a compra: ", err)
			return err
		}

		if s.Status == "Cancelada" {
			u.ErrorLogger.Printf("Essa compra n° %d, já está cancelada", s.Id)
			return fmt.Errorf("Essa compra  n° %d, já está cancelada", s.Id)
		}

		var shoppingItensGridId int

		if err := tx.QueryRow(
			ctx,
			`
				SELECT
					id
				FROM
					shopping_itens_grid
				WHERE
					shopping_id = $1
				LIMIT
					1
			`,
			c.ShoppingId,
		).Scan(
			&shoppingItensGridId,
		); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			u.ErrorLogger.Println("Erro ao conferir se possui grids na compra:", err)
			return err
		}

		if shoppingItensGridId == 0 {
			itensShoppingRows, err := tx.Query(
				ctx,
				`
					SELECT
						product_id,
						qtde_purchased
					FROM
						shopping_itens
					WHERE
						shopping_id = $1
				`,
				c.ShoppingId,
			)

			if err != nil {
				u.ErrorLogger.Println("Erro ao executar o query: ", err)
				return err
			}

			defer itensShoppingRows.Close()

			var shoppingProducts []struct {
				ProductId int
				Qtde      int
			}

			for itensShoppingRows.Next() {
				var sp struct {
					ProductId int
					Qtde      int
				}

				if err := itensShoppingRows.Scan(
					&sp.ProductId,
					&sp.Qtde,
				); err != nil {
					u.ErrorLogger.Println("Erro ao conferir os dados: ", err)
					return err
				}

				shoppingProducts = append(shoppingProducts, sp)
			}

			for _, p := range shoppingProducts {
				product := product.ProductContract{
					Id: p.ProductId,
				}

				if err := product.DiscountedQtde(ctx, tx, p.Qtde, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao discontar a qtde do item da compra:", err)
					return err
				}
			}
		} else {
			itensShoppingGridRows, err := tx.Query(
				ctx,
				`
					SELECT
						product_id,
						size_saled,
						grid_qtde
					FROM
						shopping_itens_grid
					WHERE				
						shopping_id = $1
				`,
				c.ShoppingId,
			)

			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				u.ErrorLogger.Println("Erro ao executar a query: ", err)
				return err
			}

			var productGrids []product.ProductGrids

			defer itensShoppingGridRows.Close()

			for itensShoppingGridRows.Next() {
				var productGrid product.ProductGrids

				if err := itensShoppingGridRows.Scan(
					&productGrid.ProductId,
					&productGrid.Size,
					&productGrid.GridQtde,
				); err != nil {
					u.ErrorLogger.Println("Erro ao ler os dados das grades dos produtos:", err)
					return err
				}

				productGrids = append(productGrids, productGrid)
			}

			for _, grid := range productGrids {
				productDataId := product.ProductContract{
					Id: grid.ProductId,
				}

				productGrid := product.ProductGrids{
					ProductId: productDataId.Id,
					Size:      grid.Size,
					GridQtde:  grid.GridQtde,
				}

				u.InfoLogger.Println("Vai adicionar a qtde de grade do produto (COMPRA):", productGrid)

				if err := productDataId.DiscountedQtde(ctx, tx, 0, true, &productGrid); err != nil {
					u.ErrorLogger.Println("Erro ao descontar a qtde da grade do item ao estoque:", err)
					return err
				}
			}
		}

		u.InfoLogger.Println("Vai cancelar a compra")

		if _, err = tx.Exec(
			ctx,
			`
				UPDATE
					shopping
				SET
					status = 'Cancelada'

				WHERE
					id = $1
			`,
			c.ShoppingId,
		); err != nil {
			u.ErrorLogger.Printf("Erro no update shopping para cancelado - %s", err)
			return err
		}

		payMentFormsSelect, err := tx.Query(
			ctx,
			`
				SELECT
					specie_id,
					specie,
					amount_paid
				FROM
					shopping_pay_ment
		
				WHERE
					shopping_id = $1
			`,
			s.Id,
		)

		if err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da compra - %s", err)
			return err
		}

		defer payMentFormsSelect.Close()

		var payMentFormsFromShopping []PayMentBody

		for payMentFormsSelect.Next() {
			var pf PayMentBody

			if err := payMentFormsSelect.Scan(
				&pf.SpecieId,
				&pf.Specie,
				&pf.AmountPaid,
			); err != nil {
				u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
				return err
			}

			payMentFormsFromShopping = append(payMentFormsFromShopping, pf)
		}

		customer, err := customer.Show(1)

		if err != nil {
			u.ErrorLogger.Println("Erro no select do cliente para validar a venda: ", err)
			return err
		}

		for _, pf := range payMentFormsFromShopping {
			if err := createInCashRegister(
				tx,
				pf.AmountPaid,
				0.0,
				0,
				c.ShoppingId,
				*customer,
				pf,
			); len(err) > 0 {
				u.ErrorLogger.Printf("Erro no insert de estorno no caixa compra - %s", err)
				return fmt.Errorf("Erro ao registrar o estorno no caixa.")
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		u.ErrorLogger.Printf("Erro ao comitar - %s", err)
		return err
	}

	return nil
}
