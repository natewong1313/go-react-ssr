package go_ssr

import (
	"os"
	"path/filepath"

	"github.com/natewong1313/go-react-ssr/internal/cache"
	"github.com/natewong1313/go-react-ssr/internal/typeconverter"
	"github.com/rs/zerolog"
)

type Engine struct {
	Logger       zerolog.Logger
	Config       *Config
	HotReload    *HotReload
	CacheManager *cache.Manager
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
	// Validate the config
	if err := config.validate(); err != nil {
		engine.Logger.Err(err).Msg("Failed to validate config")
		return nil, err
	}
	// Clean and make the cache directories for the type converter and tailwind
	if err := cache.SetupCacheDirectories(); err != nil {
		engine.Logger.Err(err).Msg("Failed to setup cache directories")
		return nil, err
	}
	// If using a layout css file, build it and cache it
	if config.TailwindEnabled {
		if err := engine.BuildTailwindCSSFile(); err != nil {
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
	generatedTypesPath := filepath.Join(engine.Config.FrontendSrcDir, "generated.d.ts")
	if err := typeconverter.Start(engine.Config.PropsStructsPath, generatedTypesPath); err != nil {
		engine.Logger.Err(err).Msg("Failed to init type converter")
		return nil, err
	}

	engine.Logger.Debug().Msg("Starting hot reload server")
	engine.HotReload = newHotReload(engine)
	engine.HotReload.Start()
	return engine, nil
}
