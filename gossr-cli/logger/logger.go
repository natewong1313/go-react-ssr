package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var L zerolog.Logger

// Init initializes a global logger instance
func Init() {
	L = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
}
