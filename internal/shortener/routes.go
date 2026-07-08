package shortener

import "github.com/gin-gonic/gin"

func RegisterShortenerRoutes(router *gin.Engine, handler *ShortenerHandler) {
	router.POST("/shorten", handler.ShortenURL)
	router.GET("/:short_url", handler.RedirectURL)
}
