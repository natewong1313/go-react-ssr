package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/natewong1313/go-react-ssr/internal/utils"
)

// Go SSR config
type Config struct {
	AppEnv              string `default:"development"`                       // development or production
	AssetRoute          string `default:"/assets"`                           // Route to serve static assets from
	FrontendDir         string `default:"./frontend/src"`                    // Path to frontend directory
	GeneratedTypesPath  string `default:"./frontend/src/generated/types.ts"` // Path where generated types will be written to
	PropsStructsPath    string `default:"./api/models/props.go"`             // Path to props structs file
	GlobalCSSFilePath   string `default:""`                                  // Path to global css file
	TailwindConfigPath  string `default:""`                                  // Path to Tailwind config file
	HotReloadServerPort int    `default:"3001"`                              // Port to run hot reload server on
	LayoutFile          string `default:""`                                  // Path to layout file, used for wrapping all pages
}

var C Config

// Load loads the config and runs validations
func Load(config Config) error {
	C = config
	if !checkPathExists(C.FrontendDir) {
		return fmt.Errorf("frontend dir ar %s does not exist", C.FrontendDir)
	}
	if os.Getenv("APP_ENV") != "production" && !checkPathExists(C.PropsStructsPath) {
		return fmt.Errorf("props structs path at %s does not exist", C.PropsStructsPath)
	}
	if C.GlobalCSSFilePath != "" && !checkPathExists(C.GlobalCSSFilePath) {
		return fmt.Errorf("global css file path at %s does not exist", C.GlobalCSSFilePath)
	}
	if C.TailwindConfigPath != "" {
		if C.GlobalCSSFilePath == "" {
			return errors.New("global css file path must be provided when using tailwind")
		} else if !checkPathExists(C.GlobalCSSFilePath) {
			return fmt.Errorf("global css file path at %s does not exist", C.GlobalCSSFilePath)
		} else if !checkPathExists(C.TailwindConfigPath) {
			return fmt.Errorf("tailwind config path at %s does not exist", C.TailwindConfigPath)
		}
	}
	if C.LayoutFile != "" && !checkPathExists(C.LayoutFile) && !checkPathExists(path.Join(C.FrontendDir, C.LayoutFile)) {
		return fmt.Errorf("layout file path at %s does not exist", C.LayoutFile)
	}
	if checkPathExists(path.Join(C.FrontendDir, C.LayoutFile)) {
		C.LayoutFile = path.Join(C.FrontendDir, C.LayoutFile)
	}
	C.LayoutFile = utils.GetFullFilePath(C.LayoutFile)
	return nil
}

func checkPathExists(path string) bool {
	_, err := os.Stat(utils.GetFullFilePath(path))
	return !os.IsNotExist(err)
}
