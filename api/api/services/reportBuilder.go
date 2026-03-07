package reportservices

import (
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

func CreateReport(data map[string]interface{}) (string, error) {
	u.InfoLogger.Println("InitReportProps started")
	m := CreatePDFMaroto(data)

	doc, err := m.Generate()
	if err != nil {
		u.ErrorLogger.Println("Erro ao gerar o PDF.", err)
		return "", err
	}

	dir := "./files/"

	if err := os.MkdirAll(dir, 0o755); err != nil {
		u.ErrorLogger.Println("Erro ao criar o diretório do PDF.", err)
		return "", err
	}

	filePath := filepath.Join(dir, "Relatório.pdf")
	absPath, err := filepath.Abs(filePath)

	if err != nil {
		u.ErrorLogger.Println("Erro ao conferir o caminho do PDF.", err)
		return "", err
	}

	if err := doc.Save(absPath); err != nil {
		u.ErrorLogger.Println("Erro ao salvar o PDF.", err)
		return "", err
	}

	return filePath, nil
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

	return m
}

func getTransactions(transactionType string) []core.Row {
	switch transactionType{
	case "cash-register":
		rows := []core.Row{
			row.New(5).Add(
				col.New(3),
				text.NewCol()
			)
		}
	}
}

func getPageHeader() core.Row {
	return row.New(20).Add(
		col.New(6),
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

func getBlueColor() *props.Color {
	return &props.Color{
		Red:   10,
		Green: 10,
		Blue:  150,
	}
}

func getRedColor() *props.Color {
	return &props.Color{
		Red:   150,
		Green: 10,
		Blue:  10,
	}
}
