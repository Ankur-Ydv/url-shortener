package shortener

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const shortCodeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortenerService struct {
	repo        *ShortenerRepository
	redisClient *redis.Client
}

func NewShortenerService(repo *ShortenerRepository, redisClient *redis.Client) *ShortenerService {
	return &ShortenerService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, url string, expiresAt time.Time) (string, error) {
	shortCode, err := generateShortCode(6)
	if err != nil {
		return "", err
	}

	var record Record
	record.URL = url
	record.ShortURL = shortCode
	record.ExpiresAt = expiresAt

	if err := s.repo.Save(ctx, record); err != nil {
		return "", err
	}

	return record.ShortURL, nil
}

func (s *ShortenerService) GetURL(ctx context.Context, shortUrl string) (string, error) {
	if url, err := s.redisClient.Get(ctx, shortUrl).Result(); err == nil {
		return url, nil
	}

	url, err := s.repo.Get(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	if err := s.redisClient.Set(ctx, shortUrl, url, 24*time.Hour).Err(); err != nil {
		return "", err
	}

	return url, nil
}

func (s *ShortenerService) DeleteURL(ctx context.Context, shortUrl string) error {
	if err := s.repo.Delete(ctx, shortUrl); err != nil {
		return err
	}

	if err := s.redisClient.Del(ctx, shortUrl).Err(); err != nil {
		return err
	}

	return nil
}

func (s *ShortenerService) GetShortURL(ctx context.Context, url string) (string, error) {
	if shortUrl, err := s.redisClient.Get(ctx, url).Result(); err == nil {
		return shortUrl, nil
	}

	shortUrl, err := s.repo.GetShortURL(ctx, url)
	if err != nil {
		return "", err
	}

	if err := s.redisClient.Set(ctx, url, shortUrl, 24*time.Hour).Err(); err != nil {
		return "", err
	}

	return shortUrl, nil
}

func generateShortCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}

	chars := make([]byte, length)
	alphabetLen := big.NewInt(int64(len(shortCodeAlphabet)))

	for i := range chars {
		n, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		chars[i] = shortCodeAlphabet[n.Int64()]
	}

	return string(chars), nil
}
