package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger Logger
var Logger *zerolog.Logger

func init() {
	log := zerolog.New(os.Stdout)
	Logger = &log
}
