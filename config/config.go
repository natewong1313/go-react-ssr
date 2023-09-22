package config

import (
	"errors"
	"os"
)

type Config struct {
	FrontendDir         string `default:"./frontend/src"`
	GeneratedTypesPath  string `default:"./frontend/src/generated/types.ts"`
	PropsStructsPath    string `default:"./api/models/props.go"`
	WithTailwind        bool   `default:"false"`
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
	return nil
}

func checkPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
