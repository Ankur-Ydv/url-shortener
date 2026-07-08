package main

import (
	"context"

	"github.com/Ankur-Ydv/url-shortener/internal/config"
	"github.com/Ankur-Ydv/url-shortener/internal/shortener"
	"github.com/Ankur-Ydv/url-shortener/pkg/logger"
	"github.com/Ankur-Ydv/url-shortener/pkg/postgres"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	err = logger.InitLogger(cfg.APPENV)
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(logger.GinRequestLogger())

	dsn := "postgres://" + cfg.DBUser + ":" + cfg.DBPass + "@" + cfg.DBHost + ":" + cfg.DBPort + "/postgres?sslmode=" + cfg.DBSSLMode

	dbpool, err := postgres.GetPostgresConnection(context.Background(), dsn)
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
		return
	}
	defer postgres.ClosePostgresConnection(context.Background(), dbpool)

	shortenerRepo := shortener.NewShortenerRepository(dbpool)
	shortenerService := shortener.NewShortenerService(shortenerRepo)
	shortenerHandler := shortener.NewShortenerHandler(shortenerService)

	shortener.RegisterShortenerRoutes(router, shortenerHandler)

	if err := router.Run(":8080"); err != nil {
		logger.Log.Error("failed to start server", zap.Error(err))
		return
	}
}
