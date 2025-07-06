package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	client *redis.Client
}

func New(_ context.Context, connStr string) (*Database, error) {
	opts, err := redis.ParseURL(connStr)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	return &Database{rdb}, nil
}

func (d *Database) Client() *redis.Client {
	return d.client
}

func (d *Database) Shutdown(_ context.Context) error {
	return d.Client().Close()
}
