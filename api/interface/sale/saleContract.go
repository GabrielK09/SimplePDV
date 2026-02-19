package sale

import (
	"context"
	"fmt"
	"log"
	calchelper "myApi/helpers/calc"
	"myApi/interface/cashRegister"
	"myApi/interface/product"
	saleitem "myApi/interface/saleItem"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleContract struct {
	Id        int                       `json:"id"`
	Customer  string                    `json:"customer"`
	SaleValue float64                   `json:"sale_value"`
	Status    string                    `json:"status"`
	Products  saleitem.SaleItemContract `json:"products"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
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
	}

	if len(p.Species) <= 0 {
		errorsField["species"] = "Pagamento ausente."
	}

	for idx := range p.Species {
		payMent = p.Species[idx]

		log.Println("Forma de pagamento aqui: ", payMent)

		totalPaide += payMent.AmountPaid
	}

	log.Println("totalPaide aqui: ", totalPaide)

	if payMent.Specie != "Dinheiro" && payMent.Specie != "Pix" {
		errorsField["species.specie"] = "A espécie de pagamento precisa ser Dinheiro ou Pix."
	}

	if totalPaide <= 0 {
		errorsField["amount_paid"] = "O pagamento não pode ser menor que zero."
	}

	return errorsField
}

func (s SaleContract) Validate() map[string]string {
	var subTotal float64

	for idx := range s.Products {
		i := &s.Products[idx]

		log.Println("Produto aqui: ", i)

		subTotal += calchelper.CalculateTotalSale(i.Price, i.Qtde)
	}

	errorsField := make(map[string]string)

	if err := s.Products.Validate(); len(err) > 0 {
		errorsField["products"] = fmt.Sprintf("%s", err)
	}

	log.Println("SubTotal aqui: ", subTotal)

	if subTotal <= 0 {
		errorsField["sub_total"] = "O valor da venda não pode ser zerado."
	}

	return errorsField
}

func GetAll() ([]SaleContract, error) {
	var sales []SaleContract

	query := `
		SELECT
			id,
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
		log.Println("Erro: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s SaleContract

		if err := rows.Scan(
			&s.Id,
			&s.Customer,
			&s.SaleValue,
			&s.Status,
		); err != nil {
			log.Println("Erro: ", err)
			return nil, err
		}

		sales = append(sales, s)
	}

	return sales, nil
}

func Show(id int) (*SaleContract, error) {
	query := `
		SELECT
			id,
			customer,
			sale_value,
			status
		FROM
			sales

		WHERE
			id = $1
	`

	var s SaleContract

	err := conn.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&s.Id,
		&s.Customer,
		&s.SaleValue,
		&s.Status,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *SaleContract) Create() (int, error) {
	tx, err := conn.Begin(context.Background())

	if err != nil {
		return 0, err
	}

	defer tx.Rollback(context.Background())

	querySale := `
		INSERT INTO sales
			(customer, sale_value)

		VALUES
			($1, $2)

		RETURNING 
			id
	`

	var saleId int

	err = tx.QueryRow(
		context.Background(),
		querySale,
		s.Customer,
		s.SaleValue,
	).Scan(
		&saleId,
	)

	s.Id = saleId

	if err != nil {
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

		totalSale := calchelper.CalculateTotalSale(i.Price, i.Qtde)

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
		return err
	}

	defer tx.Rollback(ctx)

	var s SaleContract

	queryForSale := `
		SELECT
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
		&s.Customer,
		&s.SaleValue,
		&s.Status,
	)

	if err != nil {
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

	log.Println("Vai fazer o insert no sale_pay_ment pelo for")

	for _, p := range payMent.Species {
		if p.AmountPaid <= 0 {
			continue
		}

		queryForPayMent := `
			INSERT INTO sale_pay_ment
				(sale_id, specie_id, specie, amount_paid)

			VALUES
				($1, $2, $3, $4)
			
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
			return err
		}

		if err := createInCashRegister(totalPaide, payMent.SaleId, s.Customer); len(err) > 0 {
			return fmt.Errorf("Erros: %s", err)
		}
	}

	log.Println("Venda está pendente, vai finalizar a venda e os itens.")

	queryForUpdateSale := `
		UPDATE 
			sales
		SET
			status = 'Concluída'		
		WHERE 
			id = $1
	`

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
		queryForUpdateSale,
		payMent.SaleId,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		queryForSaleItem,
		payMent.SaleId,
	)

	if err != nil {
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}

func createInCashRegister(inputValue float64, saleId int, customer string) map[string]string {
	log.Println("Vai inserir no caixa")
	var c cashRegister.CashRegisterContract

	if err := c.Create(inputValue, 0.0, 0, saleId, customer); len(err) > 0 {
		return err
	}

	return nil
}

func (s SaleContract) CancelSale() error {

	return nil
}
