package go_ssr

import (
	"os"

	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/natewong1313/go-react-ssr/config"
	"github.com/natewong1313/go-react-ssr/internal/hot_reload"
	"github.com/natewong1313/go-react-ssr/internal/logger"
	"github.com/natewong1313/go-react-ssr/internal/type_converter"
	"github.com/natewong1313/go-react-ssr/react_renderer"
)

func Init(optionalCfg ...config.Config) error {
	if err := godotenv.Load(); err != nil {
		logger.L.Err(err).Msg("Error loading .env file")
	}

	logger.Init()
	cfg := getConfig(optionalCfg)
	if err := defaults.Set(cfg); err != nil {
		logger.L.Err(err).Msg("Failed to set defaults")
		return err
	}

	if err := config.Load(*cfg); err != nil {
		logger.L.Err(err).Msg("Failed to load config")
		return err
	}

	react_renderer.BuildGlobalCSSFile()

	if os.Getenv("APP_ENV") == "production" {
		logger.L.Info().Msg("Running in production mode")
		return nil
	}
	logger.L.Info().Msg("Running in development mode")
	logger.L.Debug().Msg("Starting type converter")

	if err := type_converter.Init(); err != nil {
		logger.L.Err(err).Msg("Failed to init type converter")
		return err
	}

	logger.L.Debug().Msg("Starting hot reload")
	hot_reload.Init()
	return nil
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
