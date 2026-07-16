package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Ankur-Ydv/url-shortener/internal/config"
	"github.com/Ankur-Ydv/url-shortener/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	if len(os.Args) < 2 {
		logger.Log.Fatal("usage: migrate <up|down> [steps]")
	}
	cmd := os.Args[1]

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Log.Fatal("unable to initiate migrate ", zap.Error(err))
	}
	defer m.Close()

	switch cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "version":
		if v, _, er := m.Version(); err == nil {
			logger.Log.Info("check version", zap.Any("v", v))
		} else {
			err = er
		}
	case "force":
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			logger.Log.Fatal("invalid version number", zap.Any("v", v))
		}
		err = m.Force(v)
	default:
		logger.Log.Fatal("migration failed")
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Log.Error("no change")
		} else {
			logger.Log.Fatal("migration failed", zap.Error(err))
		}
	}

	logger.Log.Info("migration completed", zap.Any("cmd", cmd))
}
