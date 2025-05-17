package config

import (
	"fmt"
	"net"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const configKeyServeStatic = "serve_static"

type ServeStaticConfigRaw struct {
	Enabled    bool   `fig:"enabled"`
	Addr       string `fig:"addr"`
	BaseApiUrl string `fig:"base_api_url"`
}

type ServeStaticConfig struct {
	Enabled    bool
	Listener   net.Listener
	BaseApiUrl string
}

type ServeStaticConfiger interface {
	ServeStaticConfig() ServeStaticConfig
}

type serveStaticConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

func NewServeStaticConfiger(getter kv.Getter) ServeStaticConfiger {
	return &serveStaticConfiger{
		getter: getter,
	}
}

func (c *serveStaticConfiger) ServeStaticConfig() ServeStaticConfig {
	return c.once.Do(func() interface{} {
		var cfgRaw ServeStaticConfigRaw

		err := figure.
			Out(&cfgRaw).
			From(kv.MustGetStringMap(c.getter, configKeyServeStatic)).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out static server config: %w", err))
		}

		cfg := ServeStaticConfig{
			Enabled:    cfgRaw.Enabled,
			BaseApiUrl: cfgRaw.BaseApiUrl,
		}

		if cfg.Enabled {
			listener, err := net.Listen("tcp", cfgRaw.Addr)
			if err != nil {
				panic(fmt.Errorf("failed to configure static listener: %w", err))
			}
			cfg.Listener = listener
		}

		return cfg
	}).(ServeStaticConfig)
}
