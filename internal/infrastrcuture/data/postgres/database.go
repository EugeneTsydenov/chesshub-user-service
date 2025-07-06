package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, connStr string) (*Database, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &Database{
		pool: pool,
	}, nil
}

func (d *Database) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *Database) Shutdown(_ context.Context) error {
	d.Pool().Close()
	return nil
}
