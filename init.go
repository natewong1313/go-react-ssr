package go_ssr

import (
	"github.com/creasty/defaults"
	"github.com/natewong1313/go-react-ssr/internal/type_converter"
	"github.com/natewong1313/go-react-ssr/pkg/config"
	"github.com/natewong1313/go-react-ssr/pkg/hot_reload"
)

func Init(optionalCfg ...config.Config) {
	cfg := getConfig(optionalCfg)
	if err := defaults.Set(cfg); err != nil {
		panic(err)
	}
	config.Load(*cfg)

	err := type_converter.Init()
	if err != nil {
		panic(err)
	}

	hot_reload.Init()
}

func getConfig(optionalCfg []config.Config) (cfg *config.Config) {
	if len(optionalCfg) > 0 {
		cfg = &optionalCfg[0]
	} else {
		cfg = &config.Config{}
	}
	return cfg
}
