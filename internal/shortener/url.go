package shortener

type urlRecord struct {
	url      string `json:"url"`
	shortUrl string `json:"short_url"`
	expiry   string `json:"expiry"`
}

type request struct {
	url    string `json:"url"`
	expiry string `json:"expiry"`
}

type response struct {
	url      string `json:"url"`
	shortUrl string `json:"short_url"`
	expiry   string `json:"expiry"`
}
