package go_ssr

import (
	"fmt"
	"github.com/natewong1313/go-react-ssr/internal/utils"
	"os"
	"path"
)

type Config struct {
	AppEnv              string
	AssetRoute          string
	FrontendDir         string
	GeneratedTypesPath  string
	PropsStructsPath    string
	LayoutCSSFilePath   string
	TailwindConfigPath  string
	HotReloadServerPort int
	LayoutFile          string
}

func (c *Config) Validate() error {
	if !checkPathExists(c.FrontendDir) {
		return fmt.Errorf("frontend dir ar %s does not exist", c.FrontendDir)
	}
	c.FrontendDir = utils.GetFullFilePath(c.FrontendDir)
	c.GeneratedTypesPath = utils.GetFullFilePath(c.GeneratedTypesPath)
	fmt.Println(c.GeneratedTypesPath)
	if os.Getenv("APP_ENV") != "production" && !checkPathExists(c.PropsStructsPath) {
		return fmt.Errorf("props structs path at %s does not exist", c.PropsStructsPath)
	}
	c.PropsStructsPath = utils.GetFullFilePath(c.PropsStructsPath)
	fmt.Println(c.PropsStructsPath)
	if c.LayoutCSSFilePath != "" {
		if !checkPathExists(c.LayoutCSSFilePath) {
			return fmt.Errorf("layout css file path at %s does not exist", c.LayoutCSSFilePath)
		}
		c.LayoutCSSFilePath = path.Join(c.FrontendDir, c.LayoutCSSFilePath)
	}

	c.TailwindConfigPath = utils.GetFullFilePath(c.TailwindConfigPath)
	c.LayoutFile = path.Join(c.FrontendDir, c.LayoutFile)
	return nil
}

func checkPathExists(path string) bool {
	_, err := os.Stat(utils.GetFullFilePath(path))
	return !os.IsNotExist(err)
}
