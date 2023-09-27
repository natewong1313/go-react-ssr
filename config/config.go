package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/natewong1313/go-react-ssr/internal/utils"
)

// Go SSR config
type Config struct {
	FrontendDir         string `default:"./frontend/src"`
	GeneratedTypesPath  string `default:"./frontend/src/generated/types.ts"`
	PropsStructsPath    string `default:"./api/models/props.go"`
	GlobalCSSFilePath   string `default:""`
	TailwindConfigPath  string `default:""`
	HotReloadServerPort int    `default:"3001"`
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
	return nil
}

func checkPathExists(path string) bool {
	_, err := os.Stat(utils.GetFullFilePath(path))
	return !os.IsNotExist(err)
}
