package logger

import (
	"go.uber.org/zap"
	"log"
)

var Logger *zap.SugaredLogger

func Initialize() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()

	if err != nil {
		log.Fatal(err)
	}

	Logger = logger.Sugar()
}
