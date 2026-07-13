package shortener

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	if err := handleCurrentDomainLoopError(req.URL); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid URL: " + err.Error()})
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

	originalUrl, err := h.svc.GetURL(ctx, shortUrl)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "URL not found"})
		return
	}

	ctx.Redirect(302, originalUrl)
}

func handleCurrentDomainLoopError(url string) error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("Might be domain host itself: %v", err)
	}

	appDomain := os.Getenv("APP_DOMAIN")
	appPort := os.Getenv("APP_PORT")

	if url == fmt.Sprintf("http://%s:%s", appDomain, appPort) || url == fmt.Sprintf("http://%s:%s/", appDomain, appPort) {
		return fmt.Errorf("loop detected for URL: %s", url)
	}
	return nil
}

func (h *ShortenerHandler) DeleteURL(ctx *gin.Context) {
	shortUrl := ctx.Param("short_url")

	if err := h.svc.DeleteURL(ctx, shortUrl); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete URL"})
		return
	}

	ctx.JSON(200, gin.H{"message": "URL deleted successfully"})
}

func (h *ShortenerHandler) GetShortURL(ctx *gin.Context) {
	url := ctx.Query("url")

	shortUrl, err := h.svc.GetShortURL(ctx, url)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Short URL not found"})
		return
	}

	ctx.JSON(200, gin.H{"short_url": shortUrl})
}
