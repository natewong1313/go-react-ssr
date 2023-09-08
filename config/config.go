package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Web struct  {
        PublicDirectory string `yaml:"public-dir"`
        SrcDirectory string `yaml:"src-dir"`
        GeneratedTypesPath string `yaml:"generated-types-path"`
    } `yaml:"web"`
}

var Config config

func LoadConfig() error {
	yamlData, err := os.ReadFile("config.yaml")
    if err != nil {
        return err
    }
    var loadedConfig config
    err = yaml.Unmarshal(yamlData, &loadedConfig)
    if err != nil {
        return err
    }
	Config = loadedConfig
	return nil
}