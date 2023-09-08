package models

type Config struct {
	Web WebConfig `yaml:"web"`
}

type WebConfig struct {
	PublicDirectory string `yaml:"public-dir"`
	SrcDirectory string `yaml:"src-dir"`
	GeneratedTypesPath string `yaml:"generated-types-path"`
}