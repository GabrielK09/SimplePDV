package dashBoard

import (
	"context"
	_ "embed"
	"errors"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/report.sql
var reportSQL string

type DashBoardContract struct {
	Commission   float64 `json:"commission"`
	TotalSales   float64 `json:"total_saled"`
	BestCustomer string  `json:"best_customer"`
	AmountSaled  float64 `json:"amount_saled"`
	StartDateStr string  `json:"start_date"`
	EndDateStr   string  `json:"end_date"`

	StartDate time.Time `json:"-"`
	EndDate   time.Time `json:"-"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func ShowDashBoard(startEnd, endDate time.Time) (DashBoardContract, error) {
	var d DashBoardContract

	u.GeneralLogger.Println("Query:\n", string(reportSQL))

	if err := conn.QueryRow(
		ctx,
		string(reportSQL),
		startEnd,
		endDate,
	).Scan(
		&d.TotalSales,
		&d.Commission,
		&d.BestCustomer,
		&d.AmountSaled,
	); err != nil {
		u.ErrorLogger.Println("Erro ao ler os dados da query: ", err)

		if errors.Is(err, pgx.ErrNoRows) {
			return DashBoardContract{
				TotalSales:   0,
				Commission:   0,
				BestCustomer: "",
				AmountSaled:  0,
			}, nil
		}

		return DashBoardContract{}, err
	}

	return d, nil
}
