package shortener

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"
)

const shortCodeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortenerService struct {
	repo *ShortenerRepository
}

func NewShortenerService(repo *ShortenerRepository) *ShortenerService {
	return &ShortenerService{
		repo: repo,
	}
}

func (s *ShortenerService) ShortenURL(ctx context.Context, url string, expiry time.Duration) (string, error) {
	shortCode, err := generateShortCode(6)
	if err != nil {
		return "", err
	}

	var record Record
	record.URL = url
	record.ShortURL = shortCode
	record.ExpiresAt = time.Now().Add(expiry)

	if err := s.repo.Save(ctx, record); err != nil {
		return "", err
	}

	return record.ShortURL, nil
}

func (s *ShortenerService) GetURL(ctx context.Context, shortUrl string) (string, error) {
	url, err := s.repo.Get(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	return url, nil
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
