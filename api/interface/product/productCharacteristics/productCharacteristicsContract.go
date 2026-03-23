package productcharacteristics

import (
	"context"
	"errors"
	"fmt"
	u "myApi/helpers/logger"
	"myApi/interface/product"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Size string

const (
	PP  Size = "PP"
	P   Size = "P"
	M   Size = "M"
	G   Size = "G"
	GG  Size = "GG"
	XG  Size = "XG"
	XGG Size = "XGG"
	EG  Size = "EG"
	EGG Size = "EGG"
	O   Size = "O"
)

type ProductCharacteristicsContract struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	Size      string    `json:"size"`
	GridQtde  int       `json:"grid_qtde"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (s Size) isValidSize() bool {
	switch s {
	case PP, P, M, G, GG, XG, XGG, EG, EGG, O:
		return true

	}

	return false
}

func (p ProductCharacteristicsContract) Validate() map[string]string {
	errorsField := make(map[string]string)

	verifyProduct, err := product.Show(p.ProductId)

	if p.ProductId <= 0 {
		errorsField["product_id"] = "Identificador do produto inválido."
	}

	if err != nil {
		errorsField["database"] = fmt.Sprintf("Ocorreu um erro ao localizar o produto: %s", err)
	}

	if verifyProduct == nil {
		errorsField["database"] = fmt.Sprintf("Produto não localizado: %s", err)
	}

	size := Size(p.Size)

	if !size.isValidSize() {
		errorsField["size"] = "Tamanho incorreto."
	}

	return errorsField
}

func GetAllByProductId(productId int) ([]ProductCharacteristicsContract, error) {
	var productGrids []ProductCharacteristicsContract

	productGridsRows, err := conn.Query(
		ctx,
		`
			SELECT
				id,
				product_id,
				size,
				grid_qtde

			FROM
				product_grids

			WHERE
				product_id = $1
		`,
		productGrids,
	)

	if err != nil {
		u.ErrorLogger.Println("Erro ao executar a query: ", err)
		return []ProductCharacteristicsContract{}, err
	}

	defer productGridsRows.Close()

	for productGridsRows.Next() {
		var p ProductCharacteristicsContract

		if err := productGridsRows.Scan(
			&p.Id,
			&p.ProductId,
			&p.Size,
			&p.GridQtde,
		); err != nil {
			u.ErrorLogger.Println("Erro ao executar a leitura dos dados: ", err)
			return []ProductCharacteristicsContract{}, err
		}

		productGrids = append(productGrids, p)
	}

	return productGrids, nil
}

func (p *ProductCharacteristicsContract) Create() error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transição: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	queryInsertNewGrid := `
		INSERT INTO product_grids
			(product_id, size, grid_qtde)

		VALUES
			($1, $2, $3)
	`

	if _, err := tx.Exec(
		ctx,
		queryInsertNewGrid,
		p.ProductId,
		p.Size,
		p.GridQtde,
	); err != nil {
		u.ErrorLogger.Println("Erro ao inserir nova grade")
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao commitar: ", err)
		return err
	}

	return nil

}

func (p *ProductCharacteristicsContract) Update() error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transição: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				product_grids
			
			SET
				size,
				grid_qtde

			WHERE	
				product_id = $1
		`,
		p.Id,
	); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o update da grade: ", err)
		return err
	}

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE
				products

			SET 
				qtde = (
					SELECT
						COALESCE(SUM(grid_qtde), 0)

					FROM
						product_grids
						
					WHERE
						product_id = $1
				)

			WHERE
				id = $1
		`,
		p.Id,
	); err != nil {
		u.ErrorLogger.Println("Erro ao atualizar a nova qtde do produto da grade: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao commitar: ", err)
		return err
	}

	return nil
}

func Show(productId int) (*ProductCharacteristicsContract, error) {
	var p ProductCharacteristicsContract
	if err := conn.QueryRow(
		ctx,
		`
			SELECT
				size,
				grid_qtde

			FROM
				product_grids
				
			WHERE	
				product_id = $1
		`,
		productId,
	).Scan(
		&p.Size,
		&p.GridQtde,
	); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao fazer o select da grade: ", err)
		return nil, err
	}

	return &p, nil
}
