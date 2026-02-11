package paymentform

import (
	"context"
	"log"
	"time"

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

func (pf PayMentForms) Validate() map[string]string {
	errorsField := make(map[string]string)

	if pf.Specie == "Pix" {
		if pf.PixKey == "" {
			errorsField["pix_key"] = "Chave PIX ausente"
		}

		if len(pf.PixKey) > 255 {
			errorsField["pix_key"] = "A chave do PIX não pode ser superior a 255 caracteres."
		}
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
		log.Println("Erro: ", err)
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
			log.Println("Erro ao ler a consulta: ", err)
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
