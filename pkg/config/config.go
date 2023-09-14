package config

type Config struct {
	FrontendDir        string `default:"./frontend/src"`
	GeneratedTypesPath string `default:"./frontend/src/generated/types.ts"`
	PropsStructsPath   string `default:"./api/models/props.go"`
}
