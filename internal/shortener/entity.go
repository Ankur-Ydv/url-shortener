package shortener

import "time"

type Record struct {
	URL       string
	ShortURL  string
	ExpiresAt time.Time
}

type Request struct {
	URL    string        `json:"url"`
	Expiry time.Duration `json:"expiry"`
}

type Response struct {
	URL      string        `json:"url"`
	ShortURL string        `json:"short_url"`
	Expiry   time.Duration `json:"expiry"`
}
