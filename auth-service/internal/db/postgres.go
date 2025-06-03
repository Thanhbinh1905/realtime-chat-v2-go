package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(dns string) error {

	pool, err := pgxpool.New(context.Background(), dns)
	if err != nil {
		return err
	}

	Pool = pool
	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
