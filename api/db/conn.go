package db

import (
	"context"
	"fmt"
	u "myApi/helpers/logger"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool
var ctx = context.Background()

func Init() (*pgxpool.Pool, error) {
	u.InfoLogger.Println("Conectando ao db ...")

	godotenv.Load()

	if Pool != nil {
		return Pool, nil
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DB_URL não definida")
	}

	ctx, cancel := context.WithTimeout(ctx, 35*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		u.ErrorLogger.Println("Erro na configuração: ", err)
		return nil, err
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		u.ErrorLogger.Println("Erro na configuração: ", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		u.ErrorLogger.Println("Erro na configuração: ", err)
		return nil, err
	}

	Pool = pool
	return pool, nil
}
