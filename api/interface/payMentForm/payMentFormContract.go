package paymentform

import (
	"context"
	"errors"
	u "myApi/helpers/logger"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PayMentForms struct {
	Id        int       `json:"id"`
	Specie    string    `json:"specie"`
	PixKey    string    `json:"pix_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool
var ctx = context.Background()

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func (p PayMentForms) Validate() map[string]string {
	errorsField := make(map[string]string)

	if len(p.PixKey) > 255 {
		u.ErrorLogger.Println("A chave do PIX não pode ser superior a 255 caracteres.")
		errorsField["pix_key"] = "A chave do PIX não pode ser superior a 255 caracteres."
	}

	return errorsField
}

func GetAll() ([]PayMentForms, error) {
	var payMents []PayMentForms

	query := `
		SELECT
			id,
			specie,
			pix_key
		FROM
			pay_ment_forms
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
		var p PayMentForms

		if err := rows.Scan(
			&p.Id,
			&p.Specie,
			&p.PixKey,
		); err != nil {
			u.ErrorLogger.Println("Erro ao ler a consulta: ", err)
			return nil, err
		}

		payMents = append(payMents, p)
	}

	return payMents, nil
}

func (p *PayMentForms) Update() (PayMentForms, error) {
	query := `
		UPDATE
			pay_ment_forms

		SET
			pix_key = $1

		WHERE
			specie = 'Pix'

		RETURNING
			id,
			specie,
			pix_key
	`

	err := conn.QueryRow(
		context.Background(),
		query,
		p.PixKey,
	).Scan(
		&p.Id,
		&p.Specie,
		&p.PixKey,
	)

	if err != nil {
		return PayMentForms{}, err
	}

	return *p, nil
}

func Show() (*PayMentForms, error) {
	query := `
		SELECT
			id,
			specie, 
			pix_key

		FROM
			pay_ment_forms

		WHERE
			specie = 'Pix'
	`

	var pf PayMentForms

	err := conn.QueryRow(
		ctx,
		query,
		pf,
	).Scan(
		&pf.Id,
		&pf.Specie,
		&pf.PixKey,
	)

	if err != nil {
		return nil, err
	}

	return &pf, nil

}

func CreateDefaultPayMents() error {
	u.InfoLogger.Println("CreateDefaultPayMents started")

	var p PayMentForms

	tx, err := conn.Begin(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao iniciar a transiction: ", err)
		return err
	}

	defer tx.Rollback(ctx)

	selectQuery := `
		SELECT
			id
		FROM
			pay_ment_forms

		LIMIT 1

	`

	if err = tx.QueryRow(
		ctx,
		selectQuery,
	).Scan(&p.Id); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.ErrorLogger.Println("Erro ao conferir se os pagamentos existem: ", err)
		return err
	}

	if p.Id > 0 {
		u.InfoLogger.Println("Os pagamentos existem")
		return nil
	}

	u.GeneralLogger.Println("Não possui valores: ", p.Id)

	query := `
		INSERT INTO pay_ment_forms 
			(id, specie, pix_key)

		VALUES
			(1, 'Dinheiro', ''),
			(2, 'Pix', '')
	`

	if _, err = tx.Exec(
		ctx,
		query,
	); err != nil {
		u.ErrorLogger.Println("Erro ao fazer o insert: ", err)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		u.ErrorLogger.Println("Erro ao fazer o commit: ", err)
		return err
	}

	u.GeneralLogger.Println("Espécies padrões cadastradas com sucesso!")

	return nil
}
