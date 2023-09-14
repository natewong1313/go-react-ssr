package config

type Config struct {
	FrontendDir         string `default:"./frontend/src"`
	GeneratedTypesPath  string `default:"./frontend/src/generated/types.ts"`
	PropsStructsPath    string `default:"./api/models/props.go"`
	HotReloadServerPort int    `default:"3001"`
}

var C Config

func Load(config Config) {
	C = config
}
