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

func createInCashRegister(
	tx pgx.Tx,
	inputValue,
	outputValue float64,
	saleId,
	shoppingId int,
	customer customer.CustomerContract,
	specie PayMentBody,
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

func (p PayContract) ValidatePay(id int) map[string]string {
	errorsField := make(map[string]string)

	var totalPaide float64
	var payMent PayMentBody

	if p.ShoppingId == 0 && p.SaleId > 0 {
		saleData, err := sale.Show(p.SaleId)

		if err != nil {
			errorsField["database"] = fmt.Sprintf("Ocorreu um erro ao conferir se a venda existe: %s", err)
			return errorsField
		}

		if saleData == nil {
			errorsField["sale_id"] = fmt.Sprintf("O identificador da venda está incorreto, %s", err)
			return errorsField
		}
	}

	if p.SaleId == 0 && p.ShoppingId > 0 {
		shoppingData, err := shopping.Show(p.ShoppingId)

		if err != nil {
			errorsField["database"] = fmt.Sprintf("Ocorreu um erro ao conferir se a compra existe: %s", err)
			return errorsField
		}

		if shoppingData == nil {
			errorsField["shopping_id"] = fmt.Sprintf("O identificador da compra está incorreto, %s", err)
			return errorsField
		}
	}

	if len(p.Species) == 0 {
		errorsField["species"] = "Pagamento ausente."
		return errorsField
	}

	u.InfoLogger.Println("Formas de pagamento: ", p.Species)

	for _, payMent = range p.Species {
		if payMent.Specie != "Dinheiro" && payMent.Specie != "Pix" {
			errorsField["species.specie"] = "A espécie de pagamento precisa ser Dinheiro ou Pix."
			return errorsField
		}

		u.GeneralLogger.Println("Forma de pagamento aqui: ", payMent)

		totalPaide += payMent.AmountPaid
	}

	u.GeneralLogger.Println("totalPaide aqui: ", totalPaide)

	if totalPaide <= 0 {
		errorsField["amount_paid"] = "O pagamento não pode ser menor que zero."
		return errorsField
	}

	return errorsField
}

func PayMentShoppingOrSale(payMent PayContract) error {
	tx, err := conn.Begin(ctx)
	var totalPaide float64

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transiction no PayMentShoppingOrSale: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	for _, p := range payMent.Species {
		totalPaide += p.AmountPaid
	}

	if payMent.SaleId > 0 {
		var s sale.SaleContract

		if err = tx.QueryRow(
			ctx,
			`
				SELECT
					customer_id,
					customer,
					sale_value,
					status
				FROM
					sales
				WHERE
					id = $1
			`,
			payMent.SaleId,
		).Scan(
			&s.CustomerId,
			&s.Customer,
			&s.SaleValue,
			&s.Status,
		); err != nil {
			u.ErrorLogger.Println("Erro no select da venda no PayMentShoppingOrSale: ", err)
			return err
		}

		if s.Status == "Concluída" {
			return fmt.Errorf("Essa venda já está finalizada.")
		}

		if totalPaide < s.SaleValue {
			u.ErrorLogger.Printf("Valor informado menor do que da venda, valor da venda: %2.f, valor pago: %2.f", s.SaleValue, totalPaide)
			return fmt.Errorf("Valor informado menor do que da venda.")
		}

		u.GeneralLogger.Println("Vai fazer o insert no sale_pay_ment pelo for")

		for _, p := range payMent.Species {
			if p.AmountPaid <= 0 {
				continue
			}

			if _, err = tx.Exec(
				ctx,
				`
					INSERT INTO sale_pay_ment
						(
							sale_id, 
							specie_id, 
							specie, 
							amount_paid
						)
		
					VALUES
						(
							$1, 
							$2, 
							$3, 
							$4
						)
					
					RETURNING
						id
				`,
				payMent.SaleId,
				p.SpecieId,
				p.Specie,
				p.AmountPaid,
			); err != nil {
				u.ErrorLogger.Println("Erro no insert do sale_pay_ment no paySale: ", err)
				return err

			}

			c, err := customer.Show(s.CustomerId)

			if err != nil {
				u.ErrorLogger.Println("Erro ao pegar os dados do cliente: ", err)
				return err
			}

			if err := createInCashRegister(
				tx,
				p.AmountPaid,
				0.0,
				payMent.SaleId,
				0,
				*c,
				p,
			); len(err) > 0 {
				return fmt.Errorf("Erros: %s", err)
			}
		}

		u.GeneralLogger.Println("Venda está pendente, vai finalizar a venda e os itens.")

		if _, err = tx.Exec(
			ctx,
			`
				UPDATE 
					sales
				SET
					status = 'Concluída'		
				WHERE 
					id = $1
			`,
			payMent.SaleId,
		); err != nil {
			u.ErrorLogger.Println("Erro no update da venda para Concluída: ", err)
			return err
		}

		if _, err = tx.Exec(
			ctx,
			`
				UPDATE
					sale_itens
				SET
					status = 'Concluída'	
				WHERE 
					sale_id = $1
			`,
			payMent.SaleId,
		); err != nil {
			u.ErrorLogger.Println("Erro no update dos itens da venda para Concluída: ", err)
			return err
		}

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
			`,
			payMent.SaleId,
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
				payMent.ShoppingId,
			)

			if err != nil {
				u.ErrorLogger.Println("Erro ao pegar os itens da venda:", err)
				return err
			}

			var saleItens []sale.SaleItensContract
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

				saleItens = append(saleItens, saleItem)
			}

			for _, i := range saleItens {
				product := product.ProductContract{
					Id: i.ProductId,
				}

				u.InfoLogger.Printf("Vai descontar qtde %d ao item: %d", i.Qtde, product.Id)

				if err := product.DiscountedQtde(ctx, tx, i.Qtde, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao descontar a qtde do item ao estoque:", err)
					return err
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
						sale_id = $1
				`,
				payMent.SaleId,
			)

			if err != nil {
				u.ErrorLogger.Println("Erro ao pegar a grade dos itens da venda:", err)
				return err
			}

			var saleItensGrids []product.ProductGrids
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

				saleItensGrids = append(saleItensGrids, saleItemGrids)
			}

			for _, grid := range saleItensGrids {
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

	} // Processo para pagamento de VENDA

	if payMent.ShoppingId > 0 {
		var s shopping.ShoppingContract

		if err = tx.QueryRow(
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
			u.ErrorLogger.Println("Erro no select da compra no PayMentShoppingOrSale: ", err)
			return err
		}

		if totalPaide < s.TotalShopping {
			return fmt.Errorf("Valor informado menor do que da compra.")
		}

		for _, p := range payMent.Species {
			if p.AmountPaid <= 0 {
				continue
			}

			if _, err = tx.Exec(
				ctx,
				`
					INSERT INTO shopping_pay_ment
						(
							shopping_id, 
							specie_id, 
							specie, 
							amount_paid
						)
		
					VALUES
						(
							$1, 
							$2, 
							$3, 
							$4
						)
					
					RETURNING
						id
				`,
				payMent.ShoppingId,
				p.SpecieId,
				p.Specie,
				p.AmountPaid,
			); err != nil {
				u.ErrorLogger.Println("Erro no insert do shopping_pay_ment no paySale: ", err)
				return err

			}

			c, err := customer.Show(1)

			if err != nil {
				u.ErrorLogger.Println("Erro ao pegar os dados do cliente: ", err)
				return err
			}

			if err := createInCashRegister(
				tx,
				0.0,
				p.AmountPaid,
				0,
				payMent.ShoppingId,
				*c,
				p,
			); len(err) > 0 {
				return fmt.Errorf("Erros: %s", err)
			}
		}

		u.GeneralLogger.Println("A compra está pendente, vai finalizar a compra e os itens.")

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
			payMent.ShoppingId,
		); err != nil {
			u.ErrorLogger.Println("Erro no update da compra para Concluída: ", err)
			return err
		}

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
			payMent.ShoppingId,
		); err != nil {
			u.ErrorLogger.Println("Erro no update da compra para Concluída: ", err)
			return err
		}

		// Conferir questão de uso de grade ou não
		u.InfoLogger.Println("Concluíu a compra e os itens, vai adicionar a qtde de cada item")

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
			payMent.ShoppingId,
		).Scan(&shoppingItensGridId); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			u.ErrorLogger.Println("Erro ao conferir se existe itens com grade na compra:", err)
			return err
		}

		if shoppingItensGridId == 0 {
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
				payMent.ShoppingId,
			)

			if err != nil {
				u.ErrorLogger.Println("Erro ao pegar os itens da compra:", err)
				return err
			}

			var shoppingItens []shopping.ShoppingItenContract
			defer shoppingItensRows.Close()

			for shoppingItensRows.Next() {
				var shoppingItem shopping.ShoppingItenContract

				if err := shoppingItensRows.Scan(
					&shoppingItem.ProductId,
					&shoppingItem.QtdePurchased,
					&shoppingItem.PurchasedValue,
				); err != nil {
					u.ErrorLogger.Println("Erro ao ler os itens da compra:", err)
					return err
				}

				shoppingItens = append(shoppingItens, shoppingItem)
			}

			for _, i := range shoppingItens {
				product := product.ProductContract{
					Id: i.ProductId,
				}

				u.InfoLogger.Printf("Vai adicionar qtde %d ao item: %d", i.QtdePurchased, product.Id)

				if err := product.AddQtde(ctx, tx, i.QtdePurchased, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao adicionar a qtde do item ao estoque:", err)
					return err
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
				payMent.ShoppingId,
			)

			if err != nil {
				u.ErrorLogger.Println("Erro ao pegar a grade dos itens da compra:", err)
				return err
			}

			var shoppingItensGrids []product.ProductGrids
			defer shoppingItensGridRows.Close()

			for shoppingItensGridRows.Next() {
				var shoppingItenGrid product.ProductGrids

				if err := shoppingItensGridRows.Scan(
					&shoppingItenGrid.ProductId,
					&shoppingItenGrid.Size,
					&shoppingItenGrid.GridQtde,
				); err != nil {
					u.ErrorLogger.Println("Erro ao pegar a grade dos itens da compra:", err)
					return err
				}

				shoppingItensGrids = append(shoppingItensGrids, shoppingItenGrid)
			}

			for _, i := range shoppingItensGrids {
				productDataId := product.ProductContract{
					Id: i.ProductId,
				}

				productGrid := product.ProductGrids{
					ProductId: productDataId.Id,
					Size:      i.Size,
					GridQtde:  i.GridQtde,
				}

				u.InfoLogger.Println("Vai adicionar a qtde de grade do produto (COMPRA):", productGrid)

				if err := productDataId.AddQtde(ctx, tx, 0, true, &productGrid); err != nil {
					u.ErrorLogger.Println("Erro ao adicionar a qtde do item ao estoque:", err)
					return err
				}

				u.SuccessLoger.Println("Terminou de adicionar as qtdes das grades dos itens")
			}
		}
	} // PROCESSAR O PAGAMENTO DA COMPRA

	if err = tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro no commit do PayMentShoppingOrSale: ", err)
		return err
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
			c.SaleId,
		).Scan(
			&s.CustomerId,
			&s.SaleValue,
			&s.Status,
		); err != nil {
			u.ErrorLogger.Println("Erro ao localizar a venda: ", err)
			return err
		}

		if s.Status == "Cancelada" {
			u.ErrorLogger.Printf("Essa venda n° %d já está cancelada", s.Id)
			return fmt.Errorf("Essa venda n° %d já está cancelada", s.Id)
		}

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
				LIMIT
					1
			`,
			c.SaleId,
		).Scan(
			&saleItensGridId,
		); err != nil && !errors.Is(err, pgx.ErrNoRows) {
			u.ErrorLogger.Println("Erro ao conferir se possui grids na venda:", err)
			return err
		}

		if saleItensGridId == 0 {
			itensSaleRows, err := tx.Query(
				ctx,
				`
					SELECT
						product_id,
						qtde,
						status
					FROM
						sale_itens
					WHERE
						sale_id = $1
				`,
				c.SaleId,
			)

			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				u.ErrorLogger.Println("Erro ao executar a query: ", err)
				return err
			}

			defer itensSaleRows.Close()

			var saleProducts []struct {
				ProductId int
				Qtde      int
				Status    string
			}

			for itensSaleRows.Next() {
				var product struct {
					ProductId int
					Qtde      int
					Status    string
				}

				if err := itensSaleRows.Scan(
					&product.ProductId,
					&product.Qtde,
					&product.Status,
				); err != nil {
					u.ErrorLogger.Println("Erro ao conferir os dados: ", err)
					return err
				}

				saleProducts = append(saleProducts, product)
			}

			for _, p := range saleProducts {
				productData := product.ProductContract{
					Id: p.ProductId,
				}

				if err := productData.AddQtde(ctx, tx, p.Qtde, false, nil); err != nil {
					u.ErrorLogger.Println("Erro ao retornar a qtde do produto: ", err)
					return err
				}
			}
		} else {
			itensSaleGridRows, err := tx.Query(
				ctx,
				`
					SELECT
						product_id,
						size_saled,
						grid_qtde
					FROM
						sale_itens_grid
					WHERE				
						sale_id = $1
				`,
				c.SaleId,
			)

			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				u.ErrorLogger.Println("Erro ao executar a query: ", err)
				return err
			}

			var productGrids []product.ProductGrids

			defer itensSaleGridRows.Close()

			for itensSaleGridRows.Next() {
				var productGrid product.ProductGrids

				if err := itensSaleGridRows.Scan(
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

				if err := productDataId.AddQtde(ctx, tx, 0, true, &productGrid); err != nil {
					u.ErrorLogger.Println("Erro ao retornar a qtde das grades do produto da venda:", err)
					return err
				}
			}
		}

		if _, err = tx.Exec(
			ctx,
			`
				UPDATE
					sales
				SET
					status = 'Cancelada'
				WHERE
					id = $1
			`,
			s.Id,
		); err != nil {
			u.ErrorLogger.Printf("Erro no update sale para cancelado - %s", err)
			return err
		}

		if _, err = tx.Exec(
			ctx,
			`
				UPDATE
					sale_itens

				SET
					status = 'Cancelada'

				WHERE
					sale_id = $1
			`,
			s.Id,
		); err != nil {
			u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
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
			s.Id,
		)

		if err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
			return err
		}

		defer payMentFormsSelect.Close()

		var payMentFormsFromSale []PayMentBody

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

			payMentFormsFromSale = append(payMentFormsFromSale, pf)
		}

		customer, err := customer.Show(s.CustomerId)

		if err != nil {
			u.ErrorLogger.Println("Erro no select do cliente para validar a venda: ", err)
			return err
		}

		for _, pf := range payMentFormsFromSale {
			if err := createInCashRegister(
				tx,
				0.0,
				pf.AmountPaid,
				s.Id,
				0,
				*customer,
				pf,
			); len(err) > 0 {
				u.ErrorLogger.Printf("Erro no insert de estorno no caixa venda - %s", err)
				return fmt.Errorf("Erro ao registrar o estorno no caixa.")
			}
		}
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
