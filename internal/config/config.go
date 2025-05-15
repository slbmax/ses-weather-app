package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	WeatherAPIConfiger
	MailjetConfiger
}

func New(getter kv.Getter) *Config {
	return &Config{
		Logger:             comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Databaser:          pgdb.NewDatabaser(getter),
		Listenerer:         comfig.NewListenerer(getter),
		WeatherAPIConfiger: NewWeatherAPIConfiger(getter),
		MailjetConfiger:    NewMailjetConfiger(getter),
	}
}
