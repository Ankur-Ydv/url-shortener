package shortener

import (
	"github.com/gin-gonic/gin"
)

type ShortenerHandler struct {
	svc *ShortenerService
}

func NewShortenerHandler(svc *ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{
		svc: svc,
	}
}

func (h *ShortenerHandler) ShortenURL(ctx *gin.Context) {
	var req Request
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	shortUrl, err := h.svc.ShortenURL(ctx, req.URL, req.Expiry)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to shorten URL"})
		return
	}

	resp := Response{
		URL:      req.URL,
		ShortURL: shortUrl,
		Expiry:   req.Expiry,
	}

	ctx.JSON(200, resp)
}

func (h *ShortenerHandler) RedirectURL(ctx *gin.Context) {
	shortUrl := ctx.Param("short_url")
	// Process the redirect logic
	originalUrl, err := h.svc.GetURL(ctx, shortUrl)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "URL not found"})
		return
	}
	ctx.Redirect(302, originalUrl)
}
