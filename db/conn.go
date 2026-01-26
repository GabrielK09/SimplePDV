package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func Init() (*pgxpool.Pool, error) {
	godotenv.Load()
	if Pool != nil {
		return Pool, nil
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DB_URL não definida")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, err
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	Pool = pool
	return pool, nil
}
