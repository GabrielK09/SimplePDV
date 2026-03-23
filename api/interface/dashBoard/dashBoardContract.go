package dashBoard

import (
	"context"
	_ "embed"
	"errors"
	u "myApi/helpers/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/dashBoardQuery.sql
var reportSQL string

type DashBoardContract struct {
	TotalSales          float64 `json:"total_saled"`
	Commission          float64 `json:"commission"`
	AmountSaled         float64 `json:"amount_saled"`
	BestCustomer        string  `json:"best_customer"`
	TotalShopping       float64 `json:"total_shopping"`
	AmountShoppingItens float64 `json:"amount_shopping_itens"`
	AmountShopping      float64 `json:"amount_shopping"`
	StartDateStr        string  `json:"start_date"`
	EndDateStr          string  `json:"end_date"`
}

//go:embed sql/popularItens.sql
var popularItensReport string

type PopularItens struct {
	ProdutoId     int     `json:"produto_id"`
	Produto       string  `json:"produto"`
	ItemSaleValue float64 `json:"item_sale_value"`
	Qtde          float64 `json:"qtde"`
	PerPage       int     `json:"per_page"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (pr *PopularItens) ShowPopularItens() ([]PopularItens, error) {
	u.InfoLogger.Println("Dados de ShowPopularItens: ", pr)
	var popularItens []PopularItens

	rows, err := conn.Query(
		ctx,
		popularItensReport,
		pr.PerPage,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao realizar a query: ", err)
		return []PopularItens{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var pi PopularItens

		if err := rows.Scan(
			&pi.ProdutoId,
			&pi.Produto,
			&pi.ItemSaleValue,
			&pi.Qtde,
		); err != nil {
			u.ErrorLogger.Println("Erro ao ler os dados da query: ", err)
			return []PopularItens{}, err

		}

		popularItens = append(popularItens, pi)
	}

	return popularItens, nil
}

func (ds *DashBoardContract) ShowDashBoard() (DashBoardContract, error) {
	var d DashBoardContract

	if err := conn.QueryRow(
		ctx,
		string(reportSQL),
		ds.StartDateStr,
		ds.EndDateStr,
	).Scan(
		&d.TotalSales,
		&d.Commission,
		&d.AmountSaled,
		&d.BestCustomer,
		&d.TotalShopping,
		&d.AmountShoppingItens,
		&d.AmountShopping,
	); err != nil {
		u.ErrorLogger.Println("Erro ao ler os dados da query: ", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return DashBoardContract{
				TotalSales:          0,
				Commission:          0,
				AmountSaled:         0,
				BestCustomer:        "",
				TotalShopping:       0,
				AmountShoppingItens: 0,
				AmountShopping:      0,
			}, nil
		}

		return DashBoardContract{}, err
	}

	return d, nil
}
