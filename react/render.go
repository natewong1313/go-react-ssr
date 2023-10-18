package react

import "github.com/rs/zerolog"

type RenderConfig struct {
	File     string
	Title    string
	MetaTags map[string]string
	Props    interface{}
}

type RenderTask struct {
	Logger       zerolog.Logger
	RouteID      string
	FilePath     string
	Props        string
	RenderConfig RenderConfig
}
