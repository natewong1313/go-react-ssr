package go_ssr

import (
	"fmt"
	"os"
	"path"

	"github.com/natewong1313/go-react-ssr/internal/utils"
)

// Config is the config for starting the engine
type Config struct {
	AppEnv              string // "production" or "development"
	AssetRoute          string // The route to serve assets from, e.g. "/assets"
	FrontendDir         string // The path to the frontend folder, where your React app lives
	GeneratedTypesPath  string // The path to the generated types file
	PropsStructsPath    string // The path to the Go structs file, the structs will be generated to TS types
	LayoutFilePath      string // The path to the layout file, relative to the frontend dir
	LayoutCSSFilePath   string // The path to the layout css file, relative to the frontend dir
	TailwindConfigPath  string // The path to the tailwind config file
	HotReloadServerPort int    // The port to run the hot reload server on, 3001 by default
	JSRuntime           string // The JS runtime to use, "node" or "bun"
}

// Validate validates the config
func (c *Config) Validate() error {
	if !checkPathExists(c.FrontendDir) {
		return fmt.Errorf("frontend dir ar %s does not exist", c.FrontendDir)
	}
	if os.Getenv("APP_ENV") != "production" && !checkPathExists(c.PropsStructsPath) {
		return fmt.Errorf("props structs path at %s does not exist", c.PropsStructsPath)
	}
	if c.LayoutFilePath != "" && !checkPathExists(path.Join(c.FrontendDir, c.LayoutFilePath)) {
		return fmt.Errorf("layout css file path at %s/%s does not exist", c.FrontendDir, c.LayoutCSSFilePath)
	}
	if c.LayoutCSSFilePath != "" && !checkPathExists(path.Join(c.FrontendDir, c.LayoutCSSFilePath)) {
		return fmt.Errorf("layout css file path at %s/%s does not exist", c.FrontendDir, c.LayoutCSSFilePath)
	}
	if c.TailwindConfigPath != "" && c.LayoutCSSFilePath == "" {
		return fmt.Errorf("layout css file path must be provided when using tailwind")
	}
	if c.HotReloadServerPort == 0 {
		c.HotReloadServerPort = 3001
	}
	c.setFilePaths()
	return nil
}

// setFilePaths sets any paths in the config to their absolute paths
func (c *Config) setFilePaths() {
	c.FrontendDir = utils.GetFullFilePath(c.FrontendDir)
	c.GeneratedTypesPath = utils.GetFullFilePath(c.GeneratedTypesPath)
	c.PropsStructsPath = utils.GetFullFilePath(c.PropsStructsPath)
	if c.LayoutFilePath != "" {
		c.LayoutFilePath = path.Join(c.FrontendDir, c.LayoutFilePath)
	}
	if c.LayoutCSSFilePath != "" {
		c.LayoutCSSFilePath = path.Join(c.FrontendDir, c.LayoutCSSFilePath)
	}
	if c.TailwindConfigPath != "" {
		c.TailwindConfigPath = utils.GetFullFilePath(c.TailwindConfigPath)
	}
}

func checkPathExists(path string) bool {
	_, err := os.Stat(utils.GetFullFilePath(path))
	return !os.IsNotExist(err)
}
