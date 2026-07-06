package logger

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger(env string) {
	var err error

	if env == "development" {
		Logger, err = zap.NewDevelopment()
	} else {
		Logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}
}
