package utils

import (
	"os"
	"runtime"

	"github.com/natewong1313/go-react-ssr/gossr-cli/logger"
)

func HandleError(err error) {
	if err.Error() == "^C" {
		logger.L.Info().Msg("Goodbye 👋")
	} else {
		_, filename, line, _ := runtime.Caller(1)
		logger.L.Error().Err(err).Msgf("An error occurred in [%s:%d]", filename, line)
	}
	os.Exit(1)
}
