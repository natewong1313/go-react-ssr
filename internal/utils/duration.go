package utils

import (
	"time"

	"github.com/natewong1313/go-react-ssr/internal/logger"
)

func Timer(name string) func() {
	start := time.Now()
	return func() {
		logger.L.Debug().Msgf("%s took %v", name, time.Since(start))
	}
}
