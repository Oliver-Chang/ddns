package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger Logger
var (
	logger *zap.Logger
)

func init() {
	config := zap.NewProductionConfig()
	Logger, err := config.Build()
	Logger.Info("asd")
	if err != nil {
		fmt.Println(err)
	}

}
