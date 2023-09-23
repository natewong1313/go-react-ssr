package config

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/natewong1313/go-react-ssr/internal/utils"
)

type Config struct {
	FrontendDir         string `default:"./frontend/src"`
	GeneratedTypesPath  string `default:"./frontend/src/generated/types.ts"`
	PropsStructsPath    string `default:"./api/models/props.go"`
	GlobalCSSFilePath   string `default:""`
	TailwindConfigPath  string `default:""`
	HotReloadServerPort int    `default:"3001"`
}

var C Config

func Load(config Config) error {
	C = config
	if !checkPathExists(C.FrontendDir) {
		return errors.New("frontend dir does not exist")
	}
	if !checkPathExists(C.GeneratedTypesPath) {
		return errors.New("generated types path does not exist")
	}
	if !checkPathExists(C.PropsStructsPath) {
		return errors.New("props structs path does not exist")
	}
	if C.GlobalCSSFilePath != "" && !checkPathExists(C.GlobalCSSFilePath) {
		return errors.New("global css file path does not exist")
	}
	if C.TailwindConfigPath != "" && C.GlobalCSSFilePath == "" {
		return errors.New("global css file path must be provided when using tailwind")
	}
	if C.TailwindConfigPath != "" && !checkTailwindInstalled() {
		return errors.New("tailwind is not installed")
	}
	return nil
}

func checkPathExists(path string) bool {
	_, err := os.Stat(utils.GetFullFilePath(path))
	return !os.IsNotExist(err)
}

func checkTailwindInstalled() bool {
	cmd := exec.Command("npm", "list", "--depth=0")
	cmd.Dir = C.FrontendDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "tailwindcss")
}
