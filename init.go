package go_ssr

import (
	"github.com/natewong1313/go-react-ssr/pkg/config"
	"github.com/natewong1313/go-react-ssr/pkg/type_converter"
)

func Init(optionalCfg ...config.Config) {
	cfg := getConfig(optionalCfg)
	// err := defaults.Apply(&cfg)
	// if err != nil {
	// 	panic(err)
	// }
	err := type_converter.Init(cfg)
	if err != nil {
		panic(err)
	}
}

func getConfig(optionalCfg []config.Config) (cfg *config.Config) {
	if len(optionalCfg) > 0 {
		cfg = &optionalCfg[0]
	} else {
		cfg = &config.Config{}
	}
	return cfg
}
