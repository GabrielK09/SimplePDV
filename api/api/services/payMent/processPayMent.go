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

const (
	StatusPending   = "Pendente"
	StatusCompleted = "Concluída"
	StatusCanceled  = "Cancelada"
)

func (p PayContract) TotalPaide() float64 {
	var total float64

	for _, specie := range p.Species {
		total += specie.AmountPaid
	}

	return total
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

	totalPaide := payMent.TotalPaide()

	if totalPaide <= 0 {
		errorsField["amount_paid"] = "O pagamento não pode ser menor que zero."
		return errorsField
	}

	return errorsField
}

func createInCashRegister(ctx context.Context, tx pgx.Tx, inputValue, outputValue float64, saleId, customerID, shoppingId int, specie PayMentBody) map[string]string {
	u.InfoLogger.Println("Called createInCashRegister")
	errorsField := make(map[string]string)

	var c cashRegister.CashRegisterContract
	var customer customer.CustomerContract

	if customerID <= 0 {
		errorsField["customer_id"] = "O ID do cliente precisa ser maior do que zero."
		return errorsField
	}

	err := tx.QueryRow(
		ctx,
		`
			SELECT name FROM customers WHERE id = $1
		`,
		customerID,
	).Scan(&customer.Name)

	if err != nil {
		u.ErrorLogger.Println("Erro ao localizar os dados do cliente:", err)
		errorsField["customer"] = "Erro ao localizar os dados do cliente."
		return errorsField
	}

	if errors.Is(err, pgx.ErrNoRows) {
		errorsField["customer"] = "Cliente não localizado."
		return errorsField
	}

	c.CustomerId = customerID
	c.Customer = customer.Name

	c.SpecieId = specie.SpecieId
	c.Specie = specie.Specie

	c.SaleId = saleId
	c.ShoppingId = shoppingId

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

func PayMentShoppingOrSale(ctx context.Context, payMent PayContract) error {
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

	if s.Status == StatusCompleted {
		return fmt.Errorf("Essa venda já está finalizada.")
	}

	if s.Status == StatusCanceled {
		return fmt.Errorf("Essa venda está cancelada.")
	}

	if err := processSaleStock(ctx, tx, payMent.SaleId, true); err != nil {
		u.ErrorLogger.Println("Erro ao processar o estoque:", err)
		return err
	}

	if err := insertSalePayMents(ctx, tx, payMent.SaleId, s.CustomerId, payMent); err != nil {
		u.ErrorLogger.Println("Erro ao inserir o pagamento da venda:", err)
		return err
	}

	if err := finallySale(ctx, tx, payMent.SaleId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar a venda:", err)
		return err
	}

	if err := finallySaleItens(ctx, tx, payMent.SaleId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar os itens a venda:", err)
		return err
	}

	if err := finallySaleGridItens(ctx, tx, payMent.SaleId); err != nil {
		u.ErrorLogger.Println("Erro ao finalizar a grade dos intes a venda:", err)
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
	cmdTag, err := tx.Exec(
		ctx,
		`
			UPDATE
				sale_itens_grid
			SET
				status = 'Concluída'	
			WHERE 
				sale_id = $1
		`,
		saleID,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro no update dos itens da venda para Concluída: ", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		u.InfoLogger.Printf("Sem grades para finalizar nessa venda %d", saleID)
		return nil
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

			if err := createInCashRegister(ctx, tx, specie.AmountPaid, 0.0, customerID, saleID, 0, specie); len(err) > 0 {
				return fmt.Errorf("Erros: %s", err)
			}
		}
	}

	return nil
}

func processSaleStock(ctx context.Context, tx pgx.Tx, saleID int, isToRemove bool) error {
	saleItensRows, err := tx.Query(
		ctx,
		`
			SELECT
				product_id,
				qtde
			FROM
				sale_itens
			WHERE
				sale_id = $1
		`,
		saleID,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os itens da venda:", err)
		return err
	}

	defer saleItensRows.Close()

	var saleItens []sale.SaleItensContract

	for saleItensRows.Next() {
		var saleItem sale.SaleItensContract

		if err := saleItensRows.Scan(
			&saleItem.ProductId,
			&saleItem.Qtde,
		); err != nil {
			u.ErrorLogger.Println("Erro ao ler os itens da venda:", err)
			return err
		}

		saleItens = append(saleItens, saleItem)
	}

	for _, item := range saleItens {
		if isToRemove {
			if err := product.DiscountedQtde(ctx, tx, item.ProductId, item.Qtde, false, nil); err != nil {
				u.ErrorLogger.Println("Erro ao descontar a qtde do item ao estoque:", err)
				return err
			}
		} else {
			if err := product.AddQtde(ctx, tx, item.ProductId, item.Qtde, false, nil); err != nil {
				u.ErrorLogger.Println("Erro ao adicionar a qtde do item ao estoque:", err)
				return err
			}
		}
	}

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
				sale_id = $1
		`,
		saleID,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar a grade dos itens da venda:", err)
		return err
	}

	defer saleItensGridRows.Close()

	var saleItensGrids []product.ProductGrids

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

		saleItensGrids = append(saleItensGrids, saleItemGrids)
	}

	for _, grid := range saleItensGrids {
		if isToRemove {
			if err := product.DiscountedQtde(ctx, tx, grid.ProductId, 0, true, &grid); err != nil {
				u.ErrorLogger.Println("Erro ao descontar a qtde da grade do item ao estoque:", err)
				return err
			}
		} else {
			if err := product.AddQtde(ctx, tx, grid.ProductId, 0, true, &grid); err != nil {
				u.ErrorLogger.Println("Erro ao adicionar a qtde da grade do item ao estoque:", err)
				return err
			}
		}
	}

	return nil
}

func processShoppingPayment(ctx context.Context, tx pgx.Tx, payMent PayContract) error {
	s, err := getShoppingData(ctx, tx, payMent.ShoppingId)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da compra:", err)
		return err
	}

	if s.Status == StatusCompleted {
		return fmt.Errorf("Essa compra já está finalizada.")
	}

	if s.Status == StatusCanceled {
		return fmt.Errorf("Essa compra está cancelada.")
	}

	if err := processShoppingStock(ctx, tx, payMent.ShoppingId, true); err != nil {
		u.ErrorLogger.Println("Erro ao processar o estoque da compra:", err)
		return err
	}

	if err := insertShoppingPayMents(ctx, tx, payMent.ShoppingId, payMent); err != nil {
		u.ErrorLogger.Println("Erro ao inserir o pagamento da compra:", err)
		return err
	}

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
					qtde_purchased
				FROM
					shopping_itens
				WHERE
					shopping_id = $1
			`,
			shoppingID,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao pegar os itens da compra:", err)
			return err
		}

		defer shoppingItensRows.Close()

		var shoppingItens []shopping.ShoppingItenContract

		for shoppingItensRows.Next() {
			var shoppingItem shopping.ShoppingItenContract

			if err := shoppingItensRows.Scan(
				&shoppingItem.ProductId,
				&shoppingItem.QtdePurchased,
			); err != nil {
				u.ErrorLogger.Println("Erro ao ler os itens da compra:", err)
				return err
			}

			shoppingItens = append(shoppingItens, shoppingItem)
		}

		for _, item := range shoppingItens {
			if isToAdd {
				if err := product.AddQtde(ctx, tx, item.ProductId, item.QtdePurchased, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao adicionar a qtde do item ao estoque:", err)
					return err
				}

			} else {
				if err := product.DiscountedQtde(ctx, tx, item.ProductId, item.QtdePurchased, false, nil); err != nil {
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

		var shoppingItensGrids []product.ProductGrids

		for shoppingItensGridRows.Next() {
			var shoppingItemGrids product.ProductGrids

			if err := shoppingItensGridRows.Scan(
				&shoppingItemGrids.ProductId,
				&shoppingItemGrids.Size,
				&shoppingItemGrids.GridQtde,
			); err != nil {
				u.ErrorLogger.Println("Erro ao pegar a grade dos itens da compra:", err)
				return err
			}

			shoppingItensGrids = append(shoppingItensGrids, shoppingItemGrids)
		}

		for _, grid := range shoppingItensGrids {
			if isToAdd {
				if err := product.AddQtde(ctx, tx, grid.ProductId, 0, true, &grid); err != nil {
					u.ErrorLogger.Println("Erro ao adicionar a qtde da grade do item ao estoque:", err)
					return err
				}
			} else {
				if err := product.DiscountedQtde(ctx, tx, grid.ProductId, 0, true, &grid); err != nil {
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
	cmdTag, err := tx.Exec(
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
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao finalizar os itens da compra para Concluída: ", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		u.InfoLogger.Printf("Sem grades para finalizar nessa compra %d", shoppingID)
		return nil
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
				`,
				shoppingID,
				specie.SpecieId,
				specie.Specie,
				specie.AmountPaid,
			); err != nil {
				u.ErrorLogger.Println("Erro no insert do shopping_pay_ment: ", err)
				return err

			}

			if err := createInCashRegister(ctx, tx, 0.0, specie.AmountPaid, 0, 1, shoppingID, specie); len(err) > 0 {
				return fmt.Errorf("Erros: %s", err)
			}
		}
	}

	return nil
}

func processSaleCancel(ctx context.Context, tx pgx.Tx, saleID int) error {
	u.InfoLogger.Println("Called processSaleCancel - ID: ", saleID)

	s, err := getSaleData(ctx, tx, saleID)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da venda:", err)
		return err
	}

	switch s.Status {
	case StatusCanceled:
		u.ErrorLogger.Printf("Essa venda n° %d já está cancelada", saleID)

		return fmt.Errorf("Essa venda n° %d já está cancelada", saleID)

	case StatusPending:
		u.InfoLogger.Println("A venda está pendente")

		if err := processSaleStock(ctx, tx, saleID, false); err != nil {
			u.ErrorLogger.Println("Erro ao processar o estoque no cancelamento:", err)
			return err
		}

		if err := cancelSale(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a venda:", err)
			return err
		}

		if err := cancelSaleItens(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar os itens da venda:", err)
			return err
		}

		if err := cancelSaleGridItens(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a grade dos itens da venda:", err)
			return err
		}

		u.SuccessLoger.Println("Venda cancelada.")

	case StatusCompleted:
		u.InfoLogger.Println("A venda está concluída")

		if err := backSalePayMent(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao realizar o estorno do caixa:", err)
			return err
		}

		if err := processSaleStock(ctx, tx, saleID, false); err != nil {
			u.ErrorLogger.Println("Erro ao processar o estoque no cancelamento:", err)
			return err
		}

		if err := cancelSale(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a venda:", err)
			return err
		}

		if err := cancelSaleItens(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar os itens da venda:", err)
			return err
		}

		if err := cancelSaleGridItens(ctx, tx, saleID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a grade dos itens da venda:", err)
			return err
		}

		u.SuccessLoger.Println("Venda cancelada.")
	}

	return nil
}

func cancelSale(ctx context.Context, tx pgx.Tx, saleID int) error {
	u.InfoLogger.Printf("Called cancelSale - ID %d", saleID)

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				sales
			SET
				status = $2
			WHERE
				id = $1
		`,
		saleID,
		StatusCanceled,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update sale para cancelado - %s", err)
		return err
	}

	u.InfoLogger.Println("Finish cancelSale")

	return nil
}

func cancelSaleItens(ctx context.Context, tx pgx.Tx, saleID int) error {
	u.InfoLogger.Println("Called cancelSaleItens")

	cmdTag, err := tx.Exec(
		ctx,
		`
			UPDATE
				sale_itens
			SET
				status = $2
			WHERE
				sale_id = $1
		`,
		saleID,
		StatusCanceled,
	)

	if err != nil {
		u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		u.ErrorLogger.Printf("Itens não localizados para essa venda n° %d", saleID)
		return fmt.Errorf("Itens não localizados para essa venda n° %d", saleID)
	}

	u.InfoLogger.Println("Finish cancelSaleItens")

	return nil
}

func cancelSaleGridItens(ctx context.Context, tx pgx.Tx, saleID int) error {
	u.InfoLogger.Println("Called cancelSaleGridItens")

	cmdTag, err := tx.Exec(
		ctx,
		`
			UPDATE
				sale_itens_grid
			SET
				status = $2
			WHERE
				sale_id = $1
		`,
		saleID,
		StatusCanceled,
	)

	if err != nil {
		u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		u.ErrorLogger.Printf("Sem grade de itens para essa venda n° %d", saleID)
		return nil
	}

	u.InfoLogger.Println("Finish cancelSaleGridItens")

	return nil
}

func backSalePayMent(ctx context.Context, tx pgx.Tx, saleID int) error {
	u.InfoLogger.Println("Called backSalePayMent")

	s, err := getSaleData(ctx, tx, saleID)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da venda:", err)
		return err
	}

	payMentFormsRows, err := tx.Query(
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

	defer payMentFormsRows.Close()

	var payMentSaleForms []PayMentBody

	for payMentFormsRows.Next() {
		var pf PayMentBody

		if err := payMentFormsRows.Scan(
			&pf.SpecieId,
			&pf.Specie,
			&pf.AmountPaid,
		); err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
			return err
		}

		payMentSaleForms = append(payMentSaleForms, pf)
	}

	for _, payForm := range payMentSaleForms {
		if err := createInCashRegister(ctx, tx, 0.0, s.SaleValue, s.CustomerId, saleID, 0, payForm); len(err) > 0 {
			u.ErrorLogger.Printf("Erro no insert de estorno no caixa venda - %s", err)
			return fmt.Errorf("Erro ao registrar o estorno no caixa: %s", err)
		}
	}

	u.InfoLogger.Println("Finish backSalePayMent")

	return nil
}

func processShoppingCancel(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	u.InfoLogger.Println("Called processShoppingCancel")

	s, err := getShoppingData(ctx, tx, shoppingID)

	if err != nil {
		u.ErrorLogger.Println("Erro ao pegar os dados da compra:", err)
		return err
	}

	switch s.Status {
	case StatusCanceled:
		u.ErrorLogger.Printf("Essa compra n° %d já está cancelada", s.Id)
		return fmt.Errorf("Essa compra n° %d já está cancelada", s.Id)

	case StatusPending:
		if err := processShoppingStock(ctx, tx, shoppingID, false); err != nil {
			u.ErrorLogger.Println("Erro ao processar o estoque do cancelamentoda compra:", err)
			return err
		}

		if err := cancelShopping(ctx, tx, shoppingID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a compra:", err)
			return err
		}

		if err := cancelShoppingItens(ctx, tx, shoppingID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar os itens da compra:", err)
			return err
		}

		u.SuccessLoger.Println("Itens da compra cancelados, vai cancelar a grade dos itens")

		if err := cancelShoppingGridItens(ctx, tx, shoppingID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a grade dos itens da compra:", err)
			return err
		}

		u.SuccessLoger.Println("Compra cancelada.")

	case StatusCompleted:
		if err := processShoppingStock(ctx, tx, shoppingID, false); err != nil {
			u.ErrorLogger.Println("Erro ao processar o estoque do cancelamentoda compra:", err)
			return err
		}

		if err := backShoppingPayMent(ctx, tx, shoppingID); err != nil {
			u.ErrorLogger.Println("Erro ao realizar o estorno do caixa:", err)
			return err
		}

		if err := cancelShopping(ctx, tx, shoppingID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a venda:", err)
			return err
		}

		if err := cancelShoppingItens(ctx, tx, shoppingID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar os itens da compra:", err)
			return err
		}

		if err := cancelShoppingGridItens(ctx, tx, shoppingID); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a grade dos itens da compra:", err)
			return err
		}

		u.SuccessLoger.Println("Compra cancelada.")
	}
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

func cancelShoppingItens(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				shopping_itens
			SET
				status = $2
			WHERE
				shopping_id = $1
		`,
		shoppingID,
		StatusCanceled,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update shopping_itens para cancelado - %s", err)
		return err
	}
	return nil
}

func cancelShoppingGridItens(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				shopping_itens_grid
			SET
				status = $2
			WHERE
				shopping_id = $1
		`,
		shoppingID,
		StatusCanceled,
	); err != nil {
		u.ErrorLogger.Printf("Erro no update shopping_itens_grid para cancelado - %s", err)
		return err
	}

	return nil
}

func backShoppingPayMent(ctx context.Context, tx pgx.Tx, shoppingID int) error {
	payMentFormsRows, err := tx.Query(
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
		shoppingID,
	)

	if err != nil {
		u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
		return err
	}

	defer payMentFormsRows.Close()

	var payMentShoppingForms []PayMentBody

	for payMentFormsRows.Next() {
		var pf PayMentBody

		if err := payMentFormsRows.Scan(
			&pf.SpecieId,
			&pf.Specie,
			&pf.AmountPaid,
		); err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
			return err
		}

		payMentShoppingForms = append(payMentShoppingForms, pf)
	}

	for _, specie := range payMentShoppingForms {
		if err := createInCashRegister(ctx, tx, specie.AmountPaid, 0.0, 0, 1, shoppingID, specie); len(err) > 0 {
			u.ErrorLogger.Printf("Erro no insert de estorno no caixa venda - %s", err)
			return fmt.Errorf("Erro ao registrar o estorno no caixa.")
		}
	}

	return nil
}

func CancelSaleOrShopping(ctx context.Context, cancel CancelContract) error {
	u.InfoLogger.Println("called CancelSaleOrShopping")

	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Printf("Erro ao iniciar a transação - %s", err)
		return err
	}

	defer tx.Rollback(ctx)

	if cancel.SaleId > 0 {
		if err := processSaleCancel(ctx, tx, cancel.SaleId); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a venda:", err)
			return err
		}
	}

	if cancel.ShoppingId > 0 {
		if err := processShoppingCancel(ctx, tx, cancel.ShoppingId); err != nil {
			u.ErrorLogger.Println("Erro ao cancelar a compra:", err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao commitar:", err)
		return err
	}

	return nil
}
