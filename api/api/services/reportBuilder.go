package reportservices

import (
	"fmt"
	reportsdata "myApi/api/services/reports"
	u "myApi/helpers/logger"
	"os"
	"path/filepath"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/johnfercher/maroto/v2"
)

// func CreateReport(data map[string]interface{}) (string, error) {
func CreateReport(data map[string]interface{}) ([]byte, error) {
	var fileName string

	switch data["report_type"] {
	case "cash-register":
		fileName = "Relatório_caixa.pdf"

	case "pay-ment-forms":
		fileName = "Relatório_formas_de_pagamento.pdf"

	case "saled-itens":
		fileName = "Relatório_itens_vendidos.pdf"

	case "shopping-itens":
		fileName = "Relatório_itens_comprados.pdf"

	case "shoppings":
		fileName = "Relatório_compras.pdf"
	}

	u.InfoLogger.Println("InitReportProps started")
	m := CreatePDFMaroto(data)

	doc, err := m.Generate()

	if err != nil {
		u.ErrorLogger.Println("Erro ao gerar o PDF.", err)
		return nil, err
	}

	dir := "./files/"

	if err := os.MkdirAll(dir, 0o755); err != nil {
		u.ErrorLogger.Println("Erro ao criar o diretório do PDF.", err)
		return nil, err
	}

	filePath := filepath.Join(dir, fileName)

	u.InfoLogger.Println("filePath: ", filePath)

	if err := doc.Save(filePath); err != nil {
		u.ErrorLogger.Println("Erro ao salvar o PDF.", err)
		return nil, err
	}

	return doc.GetBytes(), nil
}

func CreatePDFMaroto(data map[string]interface{}) core.Maroto {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	darkGrayColor := getDarkGrayColor()
	mrt := maroto.New(cfg)

	m := maroto.NewMetricsDecorator(mrt)

	if err := m.RegisterHeader(getPageHeader()); err != nil {
		return nil
	}

	m.AddRows(text.NewRow(10, "Relatório", props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))

	m.AddRow(7,
		text.NewCol(3, "Movimentação", props.Text{
			Top:   1.5,
			Size:  9,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
		}),
	).WithStyle(&props.Cell{BackgroundColor: darkGrayColor})

	u.InfoLogger.Println(data)

	m.AddRows(getTransactions(data)...)

	return m
}

func getTransactions(data map[string]interface{}) []core.Row {
	var rows []core.Row

	switch data["report_type"] {
	case "cash-register":
		contents := data["data"].([]reportsdata.CashRegister)

		rows = []core.Row{
			row.New(7).Add(
				text.NewCol(3, "Descrição", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "Cliente", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(1, "Espécie", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "Valor de entrada", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "Valor de saída", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(1, "Valor total de entrada", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(1, "Valor total de saída", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			),
		}

		for i, item := range contents {
			r := row.New(4).Add(
				text.NewCol(3, item.Descricao, props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, item.Cliente, props.Text{Size: 8, Align: align.Center}),
				text.NewCol(1, item.Especie, props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, fmt.Sprintf("R$ %2.f", item.ValorEntrada), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, fmt.Sprintf("R$ %2.f", item.ValoraSaida), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(1, fmt.Sprintf("R$ %2.f", item.TotalEntrada), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(1, fmt.Sprintf("R$ %2.f", item.TotalSaida), props.Text{Size: 8, Align: align.Center}),
			)

			if i%2 == 0 {
				gray := getGrayColor()

				r.WithStyle(&props.Cell{BackgroundColor: gray})
			}

			rows = append(rows, r)
		}
	case "pay-ment-forms":
		contents := data["data"].([]reportsdata.PayMentsForms)

		rows = []core.Row{
			row.New(5).Add(
				col.New(1),
				text.NewCol(2, "Espécie", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(10, "Valor pago", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			),
		}

		for i, item := range contents {
			r := row.New(4).Add(
				text.NewCol(4, item.Especie, props.Text{Size: 8, Align: align.Center}),
				text.NewCol(8, fmt.Sprintf("R$ %2.f", item.TotalPaid), props.Text{Size: 8, Align: align.Center}),
			)

			if i%2 == 0 {
				gray := getGrayColor()

				r.WithStyle(&props.Cell{BackgroundColor: gray})
			}

			rows = append(rows, r)
		}

	case "saled-itens":
		contents := data["data"].([]reportsdata.SaledItens)

		rows = []core.Row{
			row.New(5).Add(
				text.NewCol(2, "ID da venda", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "ID do produto", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(4, "Produto", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "Valor do item", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "Qtde vendida", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			),
		}

		for i, item := range contents {
			r := row.New(4).Add(
				text.NewCol(2, fmt.Sprint(item.SaleId), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, fmt.Sprint(item.ProductId), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(4, item.Produto, props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, fmt.Sprintf("R$ %2.f", item.ItemSaleValue), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, fmt.Sprint(item.Qtde), props.Text{Size: 8, Align: align.Center}),
			)

			if i%2 == 0 {
				gray := getGrayColor()

				r.WithStyle(&props.Cell{BackgroundColor: gray})
			}

			rows = append(rows, r)
		}

	case "shopping-itens":
		contents := data["data"].([]reportsdata.ShoppingItens)

		rows = []core.Row{
			row.New(5).Add(
				text.NewCol(2, "ID da compra", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "ID do produto", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(4, "Produto", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "Valor do item", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(2, "Qtde comprada", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			),
		}

		for i, item := range contents {
			r := row.New(4).Add(
				text.NewCol(2, fmt.Sprint(item.ShoppingId), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, fmt.Sprint(item.ProductId), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(4, item.Produto, props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, fmt.Sprintf("R$ %2.f", item.PurchasedValue), props.Text{Size: 8, Align: align.Center}),
				text.NewCol(2, fmt.Sprint(item.QtdePurchased), props.Text{Size: 8, Align: align.Center}),
			)

			if i%2 == 0 {
				gray := getGrayColor()

				r.WithStyle(&props.Cell{BackgroundColor: gray})
			}

			rows = append(rows, r)
		}

	case "shoppings":
		contents := data["data"].([]reportsdata.Shoppings)

		rows = []core.Row{
			row.New(7).Add(
				text.NewCol(6, "Carga", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
				text.NewCol(6, "Valor da compra", props.Text{Size: 9, Align: align.Center, Style: fontstyle.Bold}),
			),
		}

		for i, item := range contents {
			r := row.New(4).Add(
				text.NewCol(6, item.Load, props.Text{Size: 8, Align: align.Center}),
				text.NewCol(6, fmt.Sprintf("R$ %2.f", item.TotalShopping), props.Text{Size: 8, Align: align.Center}),
			)

			if i%2 == 0 {
				gray := getGrayColor()

				r.WithStyle(&props.Cell{BackgroundColor: gray})
			}

			rows = append(rows, r)
		}
	}

	return rows
}

func getPageHeader() core.Row {
	return row.New(20).Add(
		col.New(3).Add(
			text.New("Achadinhos da Ju", props.Text{
				Size:  8,
				Align: align.Center,
				Color: getRedColor(),
			}),
		),
	)
}

func getDarkGrayColor() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func getGrayColor() *props.Color {
	return &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getRedColor() *props.Color {
	return &props.Color{
		Red:   150,
		Green: 10,
		Blue:  10,
	}
}
