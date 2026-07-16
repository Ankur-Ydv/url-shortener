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
INSERT INTO records (short_url, url, expires_at) VALUES ($1, $2, $3)
`

func (r *ShortenerRepository) Save(ctx context.Context, record Record) error {

	_, err := r.dbPool.Exec(ctx, saveUrlQuery, record.ShortURL, record.URL, record.ExpiresAt)

	return err
}

const getUrlQuery = `
SELECT url FROM records WHERE short_url = $1 AND (expires_at IS NULL OR expires_at > NOW())
`

func (r *ShortenerRepository) Get(ctx context.Context, shortUrl string) (string, error) {

	rows := r.dbPool.QueryRow(ctx, getUrlQuery, shortUrl)

	var url string
	if err := rows.Scan(&url); err != nil {
		return "", err
	}

	return url, nil
}

const deleteUrlQuery = `
DELETE FROM records WHERE short_url = $1
`

func (r *ShortenerRepository) Delete(ctx context.Context, shortUrl string) error {
	_, err := r.dbPool.Exec(ctx, deleteUrlQuery, shortUrl)

	return err
}

const getShortUrlQuery = `
SELECT short_url FROM urls WHERE original_url = $1 AND (expires_at IS NULL OR expires_at > NOW())
`

func (r *ShortenerRepository) GetShortURL(ctx context.Context, url string) (string, error) {
	rows := r.dbPool.QueryRow(ctx, getShortUrlQuery, url)

	var shortUrl string
	if err := rows.Scan(&shortUrl); err != nil {
		return "", err
	}

	return shortUrl, nil
}
