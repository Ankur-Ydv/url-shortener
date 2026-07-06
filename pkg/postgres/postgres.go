package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPostgresConnection(context context.Context, dsn string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context, dsn)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if err := dbpool.Ping(context); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return dbpool, nil
}

func ClosePostgresConnection(context context.Context, dbpool *pgxpool.Pool) error {
	dbpool.Close()

	if err := dbpool.Ping(context); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
