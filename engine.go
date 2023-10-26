package go_ssr

import (
	"github.com/natewong1313/go-react-ssr/internal/cache"
	"github.com/natewong1313/go-react-ssr/internal/typeconverter"
	"github.com/natewong1313/go-react-ssr/internal/utils"
	"github.com/rs/zerolog"
	"os"
)

type Engine struct {
	Logger                  zerolog.Logger
	Config                  *Config
	HotReload               *HotReload
	CacheManager            *cache.Manager
	CachedLayoutCSSFilePath string
}

// New creates a new gossr Engine instance
func New(config Config) (*Engine, error) {
	engine := &Engine{
		Logger:       zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger(),
		Config:       &config,
		CacheManager: cache.NewManager(),
	}
	if err := os.Setenv("APP_ENV", config.AppEnv); err != nil {
		engine.Logger.Err(err).Msg("Failed to set APP_ENV environment variable")
	}
	err := config.Validate()
	if err != nil {
		engine.Logger.Err(err).Msg("Failed to validate config")
		return nil, err
	}
	utils.CleanCacheDirectories()
	// If using a layout css file, build it and cache it
	if config.LayoutCSSFilePath != "" {
		if err = engine.BuildLayoutCSSFile(); err != nil {
			engine.Logger.Err(err).Msg("Failed to build layout css file")
			return nil, err
		}
	}

	// If running in production mode, return and don't start hot reload or type converter
	if os.Getenv("APP_ENV") == "production" {
		engine.Logger.Info().Msg("Running go-ssr in production mode")
		return engine, nil
	}
	engine.Logger.Info().Msg("Running go-ssr in development mode")
	engine.Logger.Debug().Msg("Starting type converter")
	// Start the type converter to convert Go types to Typescript types
	if err := typeconverter.Start(engine.Config.PropsStructsPath, engine.Config.GeneratedTypesPath); err != nil {
		engine.Logger.Err(err).Msg("Failed to init type converter")
		return nil, err
	}

	engine.Logger.Debug().Msg("Starting hot reload server")
	engine.HotReload = newHotReload(engine)
	engine.HotReload.Start()
	return engine, nil
}
