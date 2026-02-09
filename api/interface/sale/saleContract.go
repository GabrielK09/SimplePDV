package sale

import (
	"context"
	"fmt"
	"log"
	"myApi/interface/product"
	saleitem "myApi/interface/saleItem"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleContract struct {
	Id        int                       `json:"id"`
	Customer  string                    `json:"customer"`
	Specie    string                    `json:"specie"`
	SaleValue float64                   `json:"sale_value"`
	Status    string                    `json:"status"`
	Products  saleitem.SaleItemContract `json:"products"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}
type PaySaleContract struct {
	SaleId     int     `json:"sale_id"`
	Specie     string  `json:"specie"`
	AmountPaid float64 `json:"amount_paid"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func calculateTotalSale(saleValue float64, qtde int) float64 {
	return saleValue * float64(qtde)
}

func (p PaySaleContract) ValidatePay() map[string]string {
	errorsField := make(map[string]string)

	if p.Specie != "Dinheiro" && p.Specie != "Pix" {
		errorsField["specie"] = "A espécie de pagamento precisa ser Dinheiro ou Pix."
	}

	if p.AmountPaid <= 0 {
		errorsField["amount_paid"] = "O valor informado precisa ser maior que zero."
	}

	if _, err := Show(p.SaleId); err != nil {
		errorsField["sale_id"] = fmt.Sprintf("O identificador da venda está incorreto, %s", err)
	}

	return errorsField
}

func (s SaleContract) Validate() map[string]string {
	var subTotal float64

	for idx := range s.Products {
		i := &s.Products[idx]

		log.Println("Produto aqui: ", i)

		subTotal += calculateTotalSale(i.Price, i.Qtde)

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

func Show(id int) (*SaleContract, error) {
	query := `
		SELECT
			id,
			customer,
			specie,
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
		&s.Specie,
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
			(customer, specie, sale_value)

		VALUES
			($1, $2, $3)

		RETURNING 
			id
	`

	var saleId int

	err = tx.QueryRow(
		context.Background(),
		querySale,
		s.Customer,
		s.Specie,
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

		totalSale := calculateTotalSale(i.Price, i.Qtde)

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

func PaySale(saleId int, amountPaind float64) error {
	ctx := context.Background()
	tx, err := conn.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	var s SaleContract

	queryForSale := `
		SELECT
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
		saleId,
	).Scan(
		&s.SaleValue,
		&s.Status,
	)

	if err != nil {
		return err
	}

	if s.Status == "Concluída" {
		return fmt.Errorf("Essa venda já está finalizada.")
	}

	if amountPaind < s.SaleValue {
		return fmt.Errorf("Valor informado menor do que da venda.")
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
	_, err = tx.Query(
		ctx,
		queryForUpdateSale,
		saleId,
	)

	if err != nil {
		return err
	}

	_, err = tx.Query(
		ctx,
		queryForSaleItem,
		saleId,
	)

	if err != nil {
		return err
	}

	log.Println("Vai fazer o insert no pay_ment_forms")

	queryForPayMent := `
		INSERT INTO pay_ment_forms
			(sale_id, specie, amount_paid)

		VALUES
			($1, $2, $3)
		
		RETURNING
			id			
	`

	_, err = tx.Query(
		ctx,
		queryForPayMent,
		saleId,
		s.Specie,
		amountPaind,
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
