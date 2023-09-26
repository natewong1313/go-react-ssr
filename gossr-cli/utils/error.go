package utils

import (
	"os"

	"github.com/natewong1313/go-react-ssr/gossr-cli/logger"
)

func HandleError(err error) {
	if err.Error() == "^C" {
		logger.L.Info().Msg("Goodbye 👋")
	} else {
		logger.L.Error().Err(err).Msg("An unknown error occured")
	}
	os.Exit(1)
}
