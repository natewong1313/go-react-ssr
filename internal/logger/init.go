package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var L zerolog.Logger

func Init() {
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	L = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	// L = zap.Must(zap.NewDevelopment())
}
