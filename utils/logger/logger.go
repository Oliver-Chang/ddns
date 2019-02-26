package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger Logger
var Logger *zap.Logger

func init() {
	var err error

	cfg := zap.NewProductionConfig()
	Logger, err = cfg.Build()
	if err != nil {
		fmt.Println(err)
	}
}
