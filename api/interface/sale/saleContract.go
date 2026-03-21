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
		u.GeneralLogger.Println("Produto aqui: ", p)

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

	querySale := `
		INSERT INTO sales
			(customer_id, customer, sale_value)

		VALUES
			($1, $2, $3)

		RETURNING 
			id
	`

	var saleId int

	if err = tx.QueryRow(
		ctx,
		querySale,
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

		if err = tx.QueryRow(
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
		); err != nil {
			return 0, err
		}

		i.Name = p.Name
		i.SaleId = saleId

		if err = p.DiscountedQtde(context.Background(), tx, i.Qtde); err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		return 0, err
	}

	return saleId, err
}

func (s *SaleContract) InsertNewItens() error {
	if s.Id <= 0 {
		return nil
	}

	u.InfoLogger.Println("InsertNewItens started")

	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transição: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	queryForCheckItens := `
		SELECT
			qtde,
			sale_value,
			id
		FROM
			sale_itens
		WHERE
			sale_id = $1
			AND product_id = $2
	`

	queryForInsertNewItem := `
		INSERT INTO sale_itens
			(product_id, name, qtde, sale_value, sale_id)

		VALUES
			($1, $2, $3, $4, $5)
	`

	for _, p := range s.Products {
		var itens struct {
			qtde      int
			saleValue float64
			id        int
		}

		err := tx.QueryRow(
			ctx,
			queryForCheckItens,
			s.Id,
			p.ProductId,
		).Scan(
			&itens.qtde,
			&itens.saleValue,
			&itens.id,
		)

		if err == nil {
			u.InfoLogger.Println("Produto já existente na venda.")
			continue
		}

		// Usar p.ProductId
		if _, err := tx.Exec(
			ctx,
			`
				UPDATE 
					sale_itens
				SET	
					qtde = $3,
					sale_value = $4
				WHERE
					sale_id = $1
					AND product_id = $2
			`,
			s.Id,
			p.ProductId,
			p.Qtde,
			p.SaleValue,
		); err != nil {
			u.ErrorLogger.Println("Erro ao atualizar o item com qtde maior: ", err)

			return err
		}

		if !errors.Is(err, pgx.ErrNoRows) {
			u.ErrorLogger.Println("Erro ao conferir se o item existe: ", err)

			return err
		}

		u.InfoLogger.Println("Novos produto não existentes na venda.")

		if _, err := tx.Exec(
			ctx,
			queryForInsertNewItem,
			p.ProductId,
			p.Name,
			p.Qtde,
			p.SaleValue,
			s.Id,
		); err != nil {
			u.ErrorLogger.Println("Erro ao inserir os novos itens: ", err)

			return err
		}
	}

	queryForUpdateNewTotalSale := `
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
	`

	if _, err := tx.Exec(
		ctx,
		queryForUpdateNewTotalSale,
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
