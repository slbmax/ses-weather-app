package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const configKeyWeatherAPI = "weather_api"

type WeatherAPIConfig struct {
	APIKey string `fig:"api_key,required"`
}

type WeatherAPIConfiger interface {
	WeatherAPIConfig() WeatherAPIConfig
}

type weatherAPIConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

func NewWeatherAPIConfiger(getter kv.Getter) WeatherAPIConfiger {
	return &weatherAPIConfiger{
		getter: getter,
	}
}

func (c *weatherAPIConfiger) WeatherAPIConfig() WeatherAPIConfig {
	return c.once.Do(func() interface{} {
		var cfg WeatherAPIConfig

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, configKeyWeatherAPI)).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out weather api config: %w", err))
		}

		return cfg
	}).(WeatherAPIConfig)
}
