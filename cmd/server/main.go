package main

import (
	"context"
	"fmt"

	"github.com/Ankur-Ydv/url-shortener/internal/config"
	"github.com/Ankur-Ydv/url-shortener/internal/shortener"
	"github.com/Ankur-Ydv/url-shortener/pkg/logger"
	"github.com/Ankur-Ydv/url-shortener/pkg/postgres"
	"github.com/Ankur-Ydv/url-shortener/pkg/redis"
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

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	dbpool, err := postgres.GetPostgresConnection(context.Background(), dsn)
	if err != nil {
		logger.Log.Error("failed to connect to database", zap.Error(err))
		return
	}
	defer postgres.ClosePostgresConnection(context.Background(), dbpool)

	redisclient, err := redis.NewRedisClient(context.Background(), cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		logger.Log.Error("failed to connect to redis", zap.Error(err))
		return
	}
	defer redis.CloseRedisClient(context.Background(), redisclient)

	shortenerRepo := shortener.NewShortenerRepository(dbpool)
	shortenerService := shortener.NewShortenerService(shortenerRepo, redisclient)
	shortenerHandler := shortener.NewShortenerHandler(shortenerService)

	shortener.RegisterShortenerRoutes(router, shortenerHandler)

	if err := router.Run(":8080"); err != nil {
		logger.Log.Error("failed to start server", zap.Error(err))
		return
	}

	logger.Log.Info("server started on port 8080")
}
