package shopping

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShoppingContract struct {
	Id        int       `json:"id"`
	Load      int16     `json:"load"`
	Operation string    `json:"operation"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var conn *pgxpool.Pool

func SetConnection(db *pgxpool.Pool) {
	conn = db
}

func AggregateQtdeByName(id, newQtde int) error {
	return nil
}
