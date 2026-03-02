package sale

import (
	"context"
	_ "embed"
	"fmt"
	calchelper "myApi/helpers/calc"
	u "myApi/helpers/logger"
	"myApi/interface/cashRegister"
	"myApi/interface/customer"
	"myApi/interface/product"
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

type PayMentBody struct {
	SpecieId   int     `json:"id"`
	Specie     string  `json:"specie"`
	AmountPaid float64 `json:"amount"`
}

type PaySaleContract struct {
	SaleId  int           `json:"sale_id"`
	Species []PayMentBody `json:"species"`
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

func (p PaySaleContract) ValidatePay() map[string]string {
	var totalPaide float64
	var payMent PayMentBody

	errorsField := make(map[string]string)

	if _, err := Show(p.SaleId); err != nil {
		errorsField["sale_id"] = fmt.Sprintf("O identificador da venda está incorreto, %s", err)
		return errorsField
	}

	if len(p.Species) <= 0 {
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

func (s SaleContract) Validate() map[string]string {
	errorsField := make(map[string]string)
	var subTotal float64

	for _, p := range s.Products {
		u.GeneralLogger.Println("Produto aqui: ", p)

		subTotal += calchelper.CalculateTotalSale(p.SaleValue, p.Qtde)
	}

	if subTotal <= 0 {
		errorsField["sub_total"] = "O valor da venda não pode ser zerado."
	}

	_, err := customer.Show(s.CustomerId)

	if err != nil {
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
		var p struct {
			Id        int       `json:"id"`
			SaleId    int       `json:"sale_id"`
			ProductId int       `json:"product_id"`
			Name      string    `json:"name"`
			Qtde      int       `json:"qtde"`
			SaleValue float64   `json:"price"`
			Status    string    `json:"status"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"cpdated_at"`
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

	if s.CustomerId > 1 {
		u.InfoLogger.Println("Cliente diferente do padrão")
		otherCustomer, err := customer.Show(s.CustomerId)

		if err != nil {
			u.ErrorLogger.Println("Erro ao localizar o cliente:", err)
			return 0, err
		}

		s.Customer = otherCustomer.Name
		s.CustomerId = otherCustomer.Id
	}

	querySale := `
		INSERT INTO sales
			(customer_id, customer, sale_value)

		VALUES
			($1, $2, $3)

		RETURNING 
			id
	`

	var saleId int

	err = tx.QueryRow(
		ctx,
		querySale,
		s.CustomerId,
		s.Customer,
		s.SaleValue,
	).Scan(
		&saleId,
	)

	s.Id = saleId

	if err != nil {
		u.ErrorLogger.Println("Erro no create (sale-contract): ", err)
		return 0, err
	}

	querySaleItem := `
		INSERT INTO sale_itens
			(product_id, name, qtde, sale_value, sale_id)

		VALUES
			($1, $2, $3, $4, $5)
			
		RETURNING 
			id,
			name,
			status
	`

	queryForProduct := `
		SELECT
			id,
			name,
			qtde

		FROM	
			products
		WHERE
			id = $1

		FOR UPDATE
	`

	for idx := range s.Products {
		i := &s.Products[idx]

		var p product.ProductContract

		if err := tx.QueryRow(
			context.Background(),
			queryForProduct,
			i.ProductId,
		).Scan(
			&p.Id,
			&p.Name,
			&p.Qtde,
		); err != nil {
			return 0, err
		}

		totalSale := calchelper.CalculateTotalSale(i.SaleValue, i.Qtde)

		if err != nil {
			return 0, err
		}

		err = tx.QueryRow(
			context.Background(),
			querySaleItem,
			i.ProductId,
			p.Name,
			i.Qtde,
			totalSale,
			saleId,
		).Scan(
			&i.Id,
			&i.Name,
			&i.Status,
		)

		i.Name = p.Name
		i.SaleId = saleId

		if err != nil {
			return 0, err
		}

		if err = p.DiscountedQtde(context.Background(), tx, i.Qtde); err != nil {
			return 0, err
		}
	}

	err = tx.Commit(context.Background())

	if err != nil {
		return 0, err
	}

	return saleId, err
}

func PaySale(payMent PaySaleContract) error {
	var totalPaide float64
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transiction no paySale: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	var s SaleContract

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

	err = tx.QueryRow(
		ctx,
		queryForSale,
		payMent.SaleId,
	).Scan(
		&s.CustomerId,
		&s.Customer,
		&s.SaleValue,
		&s.Status,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro no select da venda no paySale: ", err)
		return err
	}

	if s.Status == "Concluída" {
		return fmt.Errorf("Essa venda já está finalizada.")
	}

	for _, p := range payMent.Species {
		totalPaide += p.AmountPaid
	}

	if totalPaide < s.SaleValue {
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

		_, err = tx.Exec(
			ctx,
			queryForPayMent,
			payMent.SaleId,
			p.SpecieId,
			p.Specie,
			p.AmountPaid,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro no insert do sale_pay_ment no paySale: ", err)
			return err
		}

		if err := createInCashRegister(
			tx,
			p.AmountPaid,
			0.0,
			payMent.SaleId,
			s.CustomerId,
			s.Customer,
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

	_, err = tx.Exec(
		ctx,
		queryForUpdateSale,
		payMent.SaleId,
	)

	if err != nil {
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

	_, err = tx.Exec(
		ctx,
		queryForSaleItem,
		payMent.SaleId,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro no update dos itens da venda para Concluída: ", err)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro no commit do paySale da venda: ", err)
		return err
	}

	return nil
}

func createInCashRegister(
	tx pgx.Tx,
	inputValue,
	outputValue float64,
	saleId,
	customerId int,
	customer string,
	specie PayMentBody,
) map[string]string {
	errorsField := make(map[string]string)
	var c cashRegister.CashRegisterContract

	c.SpecieId = specie.SpecieId
	c.Specie = specie.Specie

	if inputValue > 0 && outputValue > 0 {
		u.ErrorLogger.Println("Um registro no caixa não pode ter um valor de entrada e um de saída no mesmo registro.")
		errorsField["input_value"] = "Um registro no caixa não pode ter um valor de entrada no mesmo registro de uma saída."
		errorsField["output_value"] = "Um registro no caixa não pode ter um valor de saída no mesmo registro de uma entrada."
	}

	if inputValue > 0 {
		if err := c.Create(tx, inputValue, 0.0, 0, saleId, customerId, customer); len(err) > 0 {
			return err
		}
	}

	if outputValue > 0 {
		if err := c.Create(tx, 0.0, outputValue, 0, saleId, customerId, customer); len(err) > 0 {
			return err
		}
	}

	return errorsField
}

func (s *SaleContract) CancelSale() (SaleContract, error) {
	var payMentFormsFromSale []PayMentBody

	if s.Status == "Cancelado" {
		u.ErrorLogger.Printf("Essa venda n° %d já está cancelada", s.Id)
		return SaleContract{}, fmt.Errorf("Essa venda n° %d já está cancelada", s.Id)
	}

	for _, p := range s.Products {
		u.GeneralLogger.Println("Conferindo se os produtos da venda já não estão cancelados.")
		if p.Status == "Cancelado" {
			u.ErrorLogger.Printf("O prduto n° %d venda n° %d já está cancelada", p.ProductId, s.Id)
			return SaleContract{}, fmt.Errorf("O prduto n° %d venda n° %d já está cancelada", p.ProductId, s.Id)
		}
	}

	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Printf("Erro ao iniciar a transação - %s", err)
		return SaleContract{}, err
	}

	defer tx.Rollback(ctx)

	queryCancelSale := `
		UPDATE
			sales

		SET
			status = 'Cancelado'

		WHERE
			id = $1
	`

	_, err = tx.Exec(
		ctx,
		queryCancelSale,
		s.Id,
	)

	if err != nil {
		u.ErrorLogger.Printf("Erro no update sale para cancelado - %s", err)
		return SaleContract{}, err
	}

	queryCancelSaleItem := `
		UPDATE
			sale_itens

		SET
			status = 'Cancelado'

		WHERE
			sale_id = $1
	`

	_, err = tx.Exec(
		ctx,
		queryCancelSaleItem,
		s.Id,
	)

	if err != nil {
		u.ErrorLogger.Printf("Erro no update sale_itens para cancelado - %s", err)
		return SaleContract{}, err
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

	rows, err := tx.Query(
		ctx,
		queryFromPayMentsForms,
		s.Id,
	)

	if err != nil {
		u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
		return SaleContract{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var pf PayMentBody

		if err := rows.Scan(
			&pf.SpecieId,
			&pf.Specie,
			&pf.AmountPaid,
		); err != nil {
			u.ErrorLogger.Printf("Erro no select das formas de pagamento da venda - %s", err)
			return SaleContract{}, err
		}

		payMentFormsFromSale = append(payMentFormsFromSale, pf)
	}

	for _, pf := range payMentFormsFromSale {
		if err := createInCashRegister(
			tx,
			0.0,
			pf.AmountPaid,
			s.Id,
			s.CustomerId,
			s.Customer,
			pf,
		); len(err) > 0 {
			u.ErrorLogger.Printf("Erro no insert de estorno no caixa venda - %s", err)
			return SaleContract{}, fmt.Errorf("Erro ao registrar o estorno no caixa.")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		u.ErrorLogger.Printf("Erro ao comitar - %s", err)
		return SaleContract{}, err
	}

	return *s, nil
}
