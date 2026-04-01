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
	Id        int        `json:"id"`
	SaleId    int        `json:"sale_id"`
	ProductId int        `json:"product_id"`
	Size      string     `json:"size"`
	GridQtde  int        `json:"grid_qtde"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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
				grid_qtde,
				created_at,
				updated_at,
				deleted_at

			FROM
				product_grids

			WHERE
				product_id = $1 AND
				deleted_at IS NULL
		`,
		productId,
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
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		); err != nil {
			u.ErrorLogger.Println("Erro ao executar a leitura dos dados: ", err)
			return []ProductCharacteristicsContract{}, err
		}

		productGrids = append(productGrids, p)
	}

	return productGrids, nil
}

func GetAll() ([]ProductCharacteristicsContract, error) {
	var productGrids []ProductCharacteristicsContract

	productGridsRows, err := conn.Query(
		ctx,
		`
			SELECT
				id,
				product_id,
				size,
				grid_qtde,
				created_at,
				updated_at,
				deleted_at

			FROM
				product_grids

			WHERE
				deleted_at IS NULL
		`,
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
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
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

	if _, err := tx.Exec(
		ctx,
		`
			INSERT INTO product_grids
				(product_id, size, grid_qtde)

			VALUES
				($1, $2, $3)

			ON CONFLICT (size, product_id)
			DO UPDATE SET
				size = EXCLUDED.size, 
				grid_qtde = EXCLUDED.grid_qtde,
				deleted_at = NULL,
				updated_at = now()
		`,
		p.ProductId,
		p.Size,
		p.GridQtde,
	); err != nil {
		u.ErrorLogger.Println("Erro ao inserir nova grade:", err)
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
		p.ProductId,
	); err != nil {
		u.ErrorLogger.Println("Erro ao alterar a qtde do produto com base na qtde das grades:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao commitar: ", err)
		return err
	}

	return nil

}

func (p *ProductCharacteristicsContract) Update(gridId, productGridId int) error {
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
				size = $3,
				grid_qtde = $4

			WHERE	
				product_id = $1 AND
				id = $2
		`,
		productGridId,
		gridId,
		p.Size,
		p.GridQtde,
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
		productGridId,
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

func Show(productId int) (*[]ProductCharacteristicsContract, error) {
	var productCharacteristics []ProductCharacteristicsContract

	characteristicsRows, err := conn.Query(
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
				product_id = $1 AND
				deleted_at IS NULL
		`,
		productId,
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao fazer o select da grade: ", err)
		return nil, err
	}

	defer characteristicsRows.Close()

	for characteristicsRows.Next() {
		var c ProductCharacteristicsContract

		if err := characteristicsRows.Scan(
			&c.Id,
			&c.ProductId,
			&c.Size,
			&c.GridQtde,
		); err != nil {
			u.ErrorLogger.Println("Erro ao fazer a leitura dos dados: ", err)
			return nil, err
		}

		productCharacteristics = append(productCharacteristics, c)
	}

	u.InfoLogger.Println("Dados a serem retornados: ", &productCharacteristics)
	return &productCharacteristics, nil
}

func ShowById(gridId, productId int) (ProductCharacteristicsContract, error) {
	var productCharacteristic ProductCharacteristicsContract

	err := conn.QueryRow(
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
				id = $1 AND
				product_id = $2
		`,
		gridId,
		productId,
	).Scan(
		&productCharacteristic.Id,
		&productCharacteristic.ProductId,
		&productCharacteristic.Size,
		&productCharacteristic.GridQtde,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Nada foi localizado", err)
		return ProductCharacteristicsContract{}, fmt.Errorf("Nem uma grade localizada para esses IDs.")
	}

	if err != nil {
		u.ErrorLogger.Println("Erro ao fazer o select da grade: ", err)
		return ProductCharacteristicsContract{}, err
	}

	return productCharacteristic, nil
}

func Delete(id, productId int, deletedAt time.Time) error {
	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transição: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	queryUpdateGridDeletedAt := `
		UPDATE
			product_grids

		SET
			deleted_at = $1

		WHERE 
			product_id = $2 AND
			id = $3
	`

	if _, err = conn.Exec(
		ctx,
		queryUpdateGridDeletedAt,
		deletedAt,
		productId,
		id,
	); err != nil {
		u.ErrorLogger.Println("Erro ao deletar a grade do produto: ", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		u.ErrorLogger.Println("Erro ao commitar: ", err)
		return err
	}

	return nil
}

func (p *ProductCharacteristicsContract) DiscountedGridQtde(ctx context.Context, tx pgx.Tx, qtde int, size Size) error {
	if qtde <= 0 {
		return fmt.Errorf("Qtde inválida: %d", qtde)
	}

	if !size.isValidSize() {
		u.ErrorLogger.Printf("A grade: %s, é inválida", size)
		return fmt.Errorf("Grade inválida, %s", size)
	}

	if _, err := tx.Exec(
		ctx,
		`
			UPDATE	
				product_grids

			SET
				grid_qtde = grid_qtde - $2

			WHERE	
				product_id = $1 AND
				size = $3
		`,
		p.ProductId,
		qtde,
		size,
	); err != nil {
		u.ErrorLogger.Println("Erro ao alterar a qtde de sale_itens_grid:", err)
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
		p.ProductId,
	); err != nil {
		u.ErrorLogger.Println("Erro ao alterar a qtde de products:", err)
		return err
	}

	return nil
}
