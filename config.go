package go_ssr

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/natewong1313/go-react-ssr/internal/utils"
)

// Config is the config for starting the engine
type Config struct {
	AppEnv           string // "production" or "development"
	AssetRoute       string // The route where assets are served from on your server, default: /assets
	PropsStructsPath string // The path to the Go structs file, the structs will be generated to TS types
	FrontendSrcDir   string // The path to the frontend src folder, where your React app lives
	LayoutFile       string // The path to the layout file, relative to the frontend dir
	TailwindEnabled  bool
}

type FrontendConfig struct {
	Path               string // The path to the frontend folder, where your React app lives
	GeneratedTypesPath string // The path where the generated types file will be created, relative to the frontend dir
	LayoutFile         string // The path to the layout file, relative to the frontend dir
	TailwindEnabled    bool   // Whether to use tailwind or not
}

func NewDefaultConfig() Config {
	return Config{
		AppEnv:           "development",
		AssetRoute:       "/assets",
		PropsStructsPath: "./models/props.go",
		FrontendSrcDir:   "../frontend-tailwind/src",
	}
}

// WithTailwind sets the config to use tailwind
func (c Config) WithTailwind() Config {
	c.TailwindEnabled = true
	return c
}

// WithLayout sets the layout file path
func (c Config) WithLayout(layoutFileName string) Config {
	c.LayoutFile = filepath.Join(c.FrontendSrcDir, layoutFileName)
	return c
}

// Validate validates the config
func (c *Config) validate() error {
	// if !checkPathExists(c.FrontendDir) {
	// 	return fmt.Errorf("frontend dir ar %s does not exist", c.FrontendDir)
	// }
	// if os.Getenv("APP_ENV") != "production" && !checkPathExists(c.PropsStructsPath) {
	// 	return fmt.Errorf("props structs path at %s does not exist", c.PropsStructsPath)
	// }
	// if c.LayoutFilePath != "" && !checkPathExists(path.Join(c.FrontendDir, c.LayoutFilePath)) {
	// 	return fmt.Errorf("layout css file path at %s/%s does not exist", c.FrontendDir, c.LayoutFilePath)
	// }
	if c.TailwindEnabled && !checkPathExists(path.Join(c.FrontendSrcDir, "../tailwind.config.js")) {
		return fmt.Errorf("tailwind config file at %s not found", path.Join(c.FrontendSrcDir, "../tailwind.config.js"))
	}
	c.formatFilePaths()
	return nil
}

// formatFilePaths sets any paths in the config to their absolute paths, fixes an issue on windows with paths
func (c *Config) formatFilePaths() {
	c.FrontendSrcDir = utils.GetFullFilePath(c.FrontendSrcDir)
	// c.GeneratedTypesPath = utils.GetFullFilePath(c.GeneratedTypesPath)
	c.PropsStructsPath = utils.GetFullFilePath(c.PropsStructsPath)
	if c.LayoutFile != "" {
		c.LayoutFile = utils.GetFullFilePath(c.LayoutFile)
	}
}

func checkPathExists(path string) bool {
	_, err := os.Stat(utils.GetFullFilePath(path))
	return !os.IsNotExist(err)
}
