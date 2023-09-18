package go_ssr

import (
	"os"

	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/natewong1313/go-react-ssr/internal/logger"
	"github.com/natewong1313/go-react-ssr/internal/type_converter"
	"github.com/natewong1313/go-react-ssr/pkg/config"
	"github.com/natewong1313/go-react-ssr/pkg/hot_reload"
)

func Init(optionalCfg ...config.Config) {
	err := godotenv.Load()
	if err != nil {
		logger.L.Error().Err(err).Msg("Error loading .env file")
	}

	logger.Init()
	cfg := getConfig(optionalCfg)
	if err := defaults.Set(cfg); err != nil {
		logger.L.Error().Err(err).Msg("Failed to set defaults")
		return
	}
	config.Load(*cfg)

	if os.Getenv("APP_ENV") == "production" {
		logger.L.Info().Msg("Running in production mode")
		return
	}
	logger.L.Info().Msg("Running in development mode")
	logger.L.Debug().Msg("Starting type converter")
	err = type_converter.Init()
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to init type converter")
		return
	}

	logger.L.Debug().Msg("Starting hot reload")
	hot_reload.Init()
}

func getConfig(optionalCfg []config.Config) (cfg *config.Config) {
	if len(optionalCfg) > 0 {
		cfg = &optionalCfg[0]
	} else {
		logger.L.Info().Msg("No config provided, using defaults")
		cfg = &config.Config{}
	}
	return cfg
}
