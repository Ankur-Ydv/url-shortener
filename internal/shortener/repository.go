package shortener

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortenerRepository interface {
	Save(ctx context.Context, url string, shortUrl string, expiry string) error
	Get(ctx context.Context, shortUrl string) (string, error)
}

type shortenerRepository struct {
	dbPool *pgxpool.Pool
}

func NewShortenerRepository(dbPool *pgxpool.Pool) ShortenerRepository {
	return &shortenerRepository{
		dbPool: dbPool,
	}
}

func (r *shortenerRepository) Save(ctx context.Context, url string, shortUrl string, expiry string) error {
	// Implement the logic to save the URL record in the database
	return nil
}

func (r *shortenerRepository) Get(ctx context.Context, shortUrl string) (string, error) {
	// Implement the logic to retrieve the original URL from the database using the short URL
	return "", nil
}
