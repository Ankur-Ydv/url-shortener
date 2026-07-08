package shortener

import "context"

type ShortenerService interface {
	ShortenURL(ctx context.Context, url string, expiry string) (string, error)
	GetOriginalURL(ctx context.Context, shortUrl string) (string, error)
}

type shortenerService struct {
	repo ShortenerRepository
}

func NewShortenerService(repo ShortenerRepository) ShortenerService {
	return &shortenerService{
		repo: repo,
	}
}

func (s *shortenerService) ShortenURL(ctx context.Context, url string, expiry string) (string, error) {
	// Implement the logic to generate a short URL and save it in the repository
	return "", nil
}

func (s *shortenerService) GetOriginalURL(ctx context.Context, shortUrl string) (string, error) {
	// Implement the logic to retrieve the original URL from the repository using the short URL
	return "", nil
}
