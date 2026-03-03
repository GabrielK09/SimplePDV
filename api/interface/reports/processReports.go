package reports

import (
	"context"
	_ "embed"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportBody struct {
	ReportType   string `json:"report_type"`
	StartDateStr string `json:"start_date"`
	EndDateStr   string `json:"end_date"`

	StartDate time.Time `json:"-"`
	EndDate   time.Time `json:"-"`
}

const (
	cashRegister = "cash-register"
	payMentForms = "pay-ment-forms"
	saledItens   = "saled-itens"
)

//go:embed sql/totalgroupByPayMentsForms.sql
var payMentFormsReport string

type PayMentsForms struct {
	Especie   string  `json:"especie"`
	TotalPaid float64 `json:"total_paid"`
}

//go:embed sql/saledItens.sql
var saledItensReport string

type SaledItens struct {
	SaleId        int     `json:"sale_id"`
	ProductId     int     `json:"product_id"`
	Produto       string  `json:"produto"`
	ItemSaleValue float64 `json:"item_sale_value"`
	Qtde          int     `json:"qtde"`
}

//go:embed sql/cashRegister.sql
var cashRegisterReport string

type CashRegister struct {
	Descricao    string  `json:"descricao"`
	Cliente      string  `json:"cliente"`
	Especie      string  `json:"especie"`
	ValorEntrada float64 `json:"valor_entrada"`
	ValoraSaida  float64 `json:"valora_saida"`
	TotalEntrada float64 `json:"total_entrada"`
	TotalSaida   float64 `json:"total_saida"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (r *ReportBody) BuildReport() (any, error) {
	switch r.ReportType {
	case cashRegister:
		var cashRegisterData []CashRegister

		rows, err := conn.Query(
			ctx,
			cashRegisterReport,
			r.StartDate,
			r.EndDate,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var c CashRegister

			if err := rows.Scan(
				&c.Descricao,
				&c.Cliente,
				&c.Especie,
				&c.ValorEntrada,
				&c.ValoraSaida,
				&c.TotalEntrada,
				&c.TotalSaida,
			); err != nil {
				u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

				return nil, err
			}

			cashRegisterData = append(cashRegisterData, c)
		}

		return cashRegisterData, nil

	case payMentForms:
		var payMentFormsData []PayMentsForms

		rows, err := conn.Query(
			ctx,
			payMentFormsReport,
			r.StartDate,
			r.EndDate,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var p PayMentsForms

			if err := rows.Scan(
				&p.Especie,
				&p.TotalPaid,
			); err != nil {
				u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

				return nil, err
			}

			payMentFormsData = append(payMentFormsData, p)
		}

		return payMentFormsData, nil

	case saledItens:
		var saledItensData []SaledItens

		rows, err := conn.Query(
			ctx,
			saledItensReport,
			r.StartDate,
			r.EndDate,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var si SaledItens

			if err := rows.Scan(
				&si.SaleId,
				&si.ProductId,
				&si.Produto,
				&si.ItemSaleValue,
				&si.Qtde,
			); err != nil {
				u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

				return nil, err
			}

			saledItensData = append(saledItensData, si)
		}

		return saledItensData, nil
	}

	return nil, nil
}
