package processpayment

import (
	"context"
	"fmt"
	u "myApi/helpers/logger"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	productcharacteristics "myApi/interface/product/productCharacteristics"
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

		queryForSale := `
			SELECT
				customer_id,
				customer,
				sale_value,
				status
			FROM
				sales
			WHERE
				id = $1
		`

		if err = tx.QueryRow(
			ctx,
			queryForSale,
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

			queryForPayMent := `
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
			`

			if _, err = tx.Exec(
				ctx,
				queryForPayMent,
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

		queryForUpdateSale := `
			UPDATE 
				sales
			SET
				status = 'Concluída'		
			WHERE 
				id = $1
		`

		if _, err = tx.Exec(
			ctx,
			queryForUpdateSale,
			payMent.SaleId,
		); err != nil {
			u.ErrorLogger.Println("Erro no update da venda para Concluída: ", err)
			return err
		}

		queryForSaleItem := `
			UPDATE
				sale_itens
			SET
				status = 'Concluída'	
			WHERE 
				sale_id = $1
		`

		if _, err = tx.Exec(
			ctx,
			queryForSaleItem,
			payMent.SaleId,
		); err != nil {
			u.ErrorLogger.Println("Erro no update dos itens da venda para Concluída: ", err)
			return err
		}
	} // Processo para pagamento de VENDA

	if payMent.ShoppingId > 0 {
		var s shopping.ShoppingContract

		querySelectShopping := `
			SELECT
				id,
				load,
				operation,
				status
			FROM			
				shopping

			WHERE
				id = $1
		`

		if err = tx.QueryRow(
			ctx,
			querySelectShopping,
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

			queryForPayMent := `
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
			`

			if _, err = tx.Exec(
				ctx,
				queryForPayMent,
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

		queryForUpdateShopping := `
			UPDATE
				shopping
			SET
				status = 'Concluída'		
			WHERE 
				id = $1
		`

		if _, err = tx.Exec(
			ctx,
			queryForUpdateShopping,
			payMent.ShoppingId,
		); err != nil {
			u.ErrorLogger.Println("Erro no update da compra para Concluída: ", err)
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro no commit do PayMentShoppingOrSale da compra: ", err)
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

		querySelectSale := `
			SELECT
				id,
				customer_id,
				sale_value,
				status

			FROM
				sales
			
			WHERE
				id = $1
		`

		if err := tx.QueryRow(
			ctx,
			querySelectSale,
			c.SaleId,
		).Scan(
			&s.Id,
			&s.CustomerId,
			&s.SaleValue,
			&s.Status,
		); err != nil {
			u.ErrorLogger.Println("Erro ao localizar a venda: ", err)
			return err
		}

		var payMentFormsFromSale []PayMentBody

		if s.Status == "Cancelada" {
			u.ErrorLogger.Printf("Essa venda n° %d já está cancelada", s.Id)
			return fmt.Errorf("Essa venda n° %d já está cancelada", s.Id)
		}

		itensSaleRows, err := tx.Query(
			ctx,
			`
				SELECT
					p.use_grid,
					si.product_id,
					si.qtde,
					si.status

				FROM
					sale_itens si

				INNER JOIN
					products p on p.id = si.product_id

				WHERE
					si.sale_id = $1
			`,
			c.SaleId,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao executar o query: ", err)
			return err
		}

		defer itensSaleRows.Close()

		var saleProducts []struct {
			UseGrid   bool
			ProductId int
			Qtde      int
			Status    string
		}

		for itensSaleRows.Next() {
			var p struct {
				UseGrid   bool
				ProductId int
				Qtde      int
				Status    string
			}

			if err := itensSaleRows.Scan(
				&p.UseGrid,
				&p.ProductId,
				&p.Qtde,
				&p.Status,
			); err != nil {
				u.ErrorLogger.Println("Erro ao conferir os dados: ", err)
				return err
			}

			saleProducts = append(saleProducts, p)
		}

		for _, p := range saleProducts {
			u.InfoLogger.Println("Vai fazer o update dos item: ", p)

			if p.UseGrid {
				u.InfoLogger.Println("Produto da venda possuia grade informada, vai retornar a qtde da grade: ", p)
				u.InfoLogger.Println("Vai localizar qual grade e qtde foi vendida.")

				var productGrid struct {
					SizeSaled productcharacteristics.Size
					GridQtde  int
				}

				if err := tx.QueryRow(
					ctx,
					`
						SELECT
							size_saled,
							grid_qtde										

						FROM
							sale_itens_grid

						WHERE
							product_id = $1 AND
							sale_id = $2
					`,
					p.ProductId,
					c.SaleId,
				).Scan(
					&productGrid.SizeSaled,
					&productGrid.GridQtde,
				); err != nil {
					u.ErrorLogger.Println("Erro ao retornar a qtde dos produtos da venda para o estoque: ", err)
					return err
				}

				u.InfoLogger.Println("Vai fazer o update na product_grids.")
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							product_grids

						SET
							grid_qtde = grid_qtde + $1

						WHERE
							product_id = $2 AND
							size = $3
					`,
					productGrid.GridQtde,
					p.ProductId,
					productGrid.SizeSaled,
				); err != nil {
					u.ErrorLogger.Println("Erro ao retornar a qtde dos produtos da venda para o estoque: ", err)
					return err
				}

				u.InfoLogger.Println("Vai fazer o update na products para a qtde total das grades .")
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
									product_id = $1 AND
									use_grid = true
							)

						WHERE
							id = $1
					`,
					p.ProductId,
				); err != nil {
					u.ErrorLogger.Println("Erro ao retornar a qtde dos produtos da venda para o estoque: ", err)
					return err
				}

			} else {
				if _, err := tx.Exec(
					ctx,
					`
						UPDATE
							products

						SET
							qtde = qtde + $2

						WHERE
							id = $1
					`,
					p.ProductId,
					p.Qtde,
				); err != nil {
					u.ErrorLogger.Println("Erro ao retornar a qtde dos produtos da venda para o estoque: ", err)
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

		queryCancelSaleItem := `
			UPDATE
				sale_itens

			SET
				status = 'Cancelada'

			WHERE
				sale_id = $1
		`

		if _, err = tx.Exec(
			ctx,
			queryCancelSaleItem,
			s.Id,
		); err != nil {
			u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
			return err
		}

		queryFromPayMentsForms := `
			SELECT
				specie_id,
				specie,
				amount_paid
			FROM
				sale_pay_ment
	
			WHERE
				sale_id = $1
		`

		payMentFormsSelect, err := tx.Query(
			ctx,
			queryFromPayMentsForms,
			s.Id,
		)

		if err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
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
		var s shopping.ShoppingContract

		querySelectShopping := `
			SELECT
				id,
				load,
				total_shopping,
				status

			FROM
				shopping
			
			WHERE
				id = $1
		`

		if err := tx.QueryRow(
			ctx,
			querySelectShopping,
			c.ShoppingId,
		).Scan(
			&s.Id,
			&s.Load,
			&s.TotalShopping,
			&s.Status,
		); err != nil {
			u.ErrorLogger.Println("Erro ao localizar a compra: ", err)
			return err
		}

		var payMentFormsFromSale []PayMentBody

		if s.Status == "Cancelado" {
			u.ErrorLogger.Printf("Essa compra n° %d, lote n° %d, já está cancelada", s.Id, s.Load)
			return fmt.Errorf("Essa compra  n° %d, lote n° %d, já está cancelada", s.Id, s.Load)
		}

		querySelectItensShopping := `
			SELECT
				product_id,
				qtde_purchased
			FROM
				shopping_itens
			WHERE
				id = $1
		`

		itensShoppingRows, err := tx.Query(
			ctx,
			querySelectItensShopping,
			s.Id,
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
			queryUpdateDiscountQtdeStock := `
				UPDATE
					products

				SET
					qtde = qtde - $2

				WHERE
					id = $1
			`

			if _, err := tx.Exec(
				ctx,
				queryUpdateDiscountQtdeStock,
				p.ProductId,
				p.Qtde,
			); err != nil {
				u.ErrorLogger.Println("Erro ao retornar a qtde dos produtos da venda para o estoque: ", err)
				return err

			}
		}

		queryCancelShopping := `
			UPDATE
				shopping

			SET
				status = 'Cancelada'

			WHERE
				id = $1
		`

		if _, err = tx.Exec(
			ctx,
			queryCancelShopping,
			s.Id,
		); err != nil {
			u.ErrorLogger.Printf("Erro no update shopping para cancelado - %s", err)
			return err
		}

		queryFromPayMentsForms := `
			SELECT
				specie_id,
				specie,
				amount_paid
			FROM
				shopping_pay_ment
	
			WHERE
				shopping_id = $1
		`

		payMentFormsSelect, err := tx.Query(
			ctx,
			queryFromPayMentsForms,
			s.Id,
		)

		if err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da compra - %s", err)
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

			payMentFormsFromSale = append(payMentFormsFromSale, pf)
		}

		customer, err := customer.Show(1)

		if err != nil {
			u.ErrorLogger.Println("Erro no select do cliente para validar a venda: ", err)
			return err
		}

		for _, pf := range payMentFormsFromSale {
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
