package reportsdata

import (
	"context"
	_ "embed"
	u "myApi/helpers/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportBody struct {
	ReportType   string `json:"report_type"`
	StartDateStr string `json:"start_date"`
	EndDateStr   string `json:"end_date"`
}

const (
	cashRegister  = "cash-register"
	payMentForms  = "pay-ment-forms"
	saledItens    = "saled-itens"
	shoppingItens = "shopping-itens"
	shoppings     = "shoppings"
)

//go:embed sql/payMentForms/totalgroupByPayMentsForms.sql
var payMentFormsReport string

type PayMentsForms struct {
	Especie   string  `json:"especie"`
	TotalPaid float64 `json:"total_paid"`
}

//go:embed sql/itens/sale/saledItens.sql
var saledItensReport string

type SaledItens struct {
	SaleId        int     `json:"sale_id"`
	ProductId     int     `json:"product_id"`
	Produto       string  `json:"produto"`
	ItemSaleValue float64 `json:"item_sale_value"`
	Qtde          int     `json:"qtde"`
}

//go:embed sql/cashRegister/cashRegister.sql
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

//go:embed sql/shopping/shoppings.sql
var shoppingsReport string

type Shoppings struct {
	Load          string  `json:"load"`
	TotalShopping float64 `json:"total_shopping"`
}

//go:embed sql/itens/shopping/shoppingItens.sql
var shoppingItensReport string

type ShoppingItens struct {
	ShoppingId     int     `json:"shopping_id"`
	ProductId      int     `json:"product_id"`
	Produto        string  `json:"produto"`
	PurchasedValue float64 `json:"purchased_value"`
	QtdePurchased  int     `json:"qtde_purchased"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (r *ReportBody) BuildDataReport() (map[string]interface{}, error) {
	dataForReturn := make(map[string]interface{})
	dataForReturn["report_type"] = r.ReportType

	switch r.ReportType {
	case cashRegister:
		var cashRegisterData []CashRegister

		rows, err := conn.Query(
			ctx,
			cashRegisterReport,
			r.StartDateStr,
			r.EndDateStr,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao executar a query dos dados: ", err)

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

		dataForReturn["data"] = cashRegisterData

	case payMentForms:
		var payMentFormsData []PayMentsForms

		rows, err := conn.Query(
			ctx,
			payMentFormsReport,
			r.StartDateStr,
			r.EndDateStr,
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

		dataForReturn["data"] = payMentFormsData

	case saledItens:
		var saledItensData []SaledItens

		rows, err := conn.Query(
			ctx,
			saledItensReport,
			r.StartDateStr,
			r.EndDateStr,
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

		dataForReturn["data"] = saledItensData

	case shoppingItens:
		var shoppingItens []ShoppingItens

		rows, err := conn.Query(
			ctx,
			shoppingItensReport,
			r.StartDateStr,
			r.EndDateStr,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var si ShoppingItens

			if err := rows.Scan(
				&si.ShoppingId,
				&si.ProductId,
				&si.Produto,
				&si.PurchasedValue,
				&si.QtdePurchased,
			); err != nil {
				u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

				return nil, err
			}

			shoppingItens = append(shoppingItens, si)
		}

		dataForReturn["data"] = shoppingItens

	case shoppings:
		var shoppinsData []Shoppings

		rows, err := conn.Query(
			ctx,
			shoppingsReport,
			r.StartDateStr,
			r.EndDateStr,
		)

		if err != nil {
			u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var si Shoppings

			if err := rows.Scan(
				&si.Load,
				&si.TotalShopping,
			); err != nil {
				u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)

				return nil, err
			}

			shoppinsData = append(shoppinsData, si)
		}

		dataForReturn["data"] = shoppinsData
	}

	u.InfoLogger.Println("Tamanho do data:", len(dataForReturn))

	return dataForReturn, nil
}
