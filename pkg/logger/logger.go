package logger

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(env string) error {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true

	if env == "development" {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	var err error
	Log, err = cfg.Build()
	if err != nil {
		return err
	}

	return nil
}

func Default() *zap.Logger {
	if Log == nil {
		if err := InitLogger("production"); err != nil {
			panic(err)
		}
	}
	return Log
}

func GinRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		start := time.Now()
		requestLogger := Default().With(
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		c.Set("request_id", requestID)
		c.Set("logger", requestLogger)

		requestLogger.Info("request started")
		c.Next()

		latency := time.Since(start)
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("route", c.FullPath()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		requestLogger.Info("request completed", fields...)
	}
}

func generateRequestID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(buf)
}
