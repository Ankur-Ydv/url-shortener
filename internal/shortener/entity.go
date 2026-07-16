package shortener

import "time"

type Record struct {
	URL       string
	ShortURL  string
	ExpiresAt time.Time
}

type Request struct {
	URL       string    `json:"url" binding:"required"`
	ExpiresAt time.Time `json:"expires_at" binding:"required"`
}

type Response struct {
	URL       string    `json:"url"`
	ShortURL  string    `json:"short_url"`
	ExpiresAt time.Time `json:"expires_at"`
}
