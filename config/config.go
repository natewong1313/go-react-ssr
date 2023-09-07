package config

import (
	"gossr/models"
	"os"

	"gopkg.in/yaml.v3"
)

var CONFIG_PATH = "config.yaml"
var Config models.Config

func LoadConfig() error {
	yamlData, err := os.ReadFile(CONFIG_PATH)
    if err != nil {
        return err
    }
    var loadedConfig models.Config
    err = yaml.Unmarshal(yamlData, &loadedConfig)
    if err != nil {
        return err
    }
	Config = loadedConfig
	return nil
}