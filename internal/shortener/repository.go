package shortener

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortenerRepository struct {
	dbPool *pgxpool.Pool
}

func NewShortenerRepository(dbPool *pgxpool.Pool) *ShortenerRepository {
	return &ShortenerRepository{
		dbPool: dbPool,
	}
}

const saveUrlQuery = `
INSERT INTO urls (short_url, original_url, expiry) VALUES ($1, $2, $3)
`

func (r *ShortenerRepository) Save(ctx context.Context, record Record) error {

	_, err := r.dbPool.Exec(ctx, saveUrlQuery, record.ShortURL, record.URL, record.ExpiresAt)

	return err
}

const getUrlQuery = `
SELECT original_url FROM urls WHERE short_url = $1 AND (expiry IS NULL OR expiry > NOW())
`

func (r *ShortenerRepository) Get(ctx context.Context, shortUrl string) (string, error) {

	rows := r.dbPool.QueryRow(ctx, getUrlQuery, shortUrl)

	var url string
	if err := rows.Scan(&url); err != nil {
		return "", err
	}

	return url, nil
}
