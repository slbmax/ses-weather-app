package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const configKeyMailjet = "mailjet"
const defaultFromName = "Weather App"

type MailjetConfig struct {
	ApiKey    string `fig:"api_key,required"`
	SecretKey string `fig:"secret_key,required"`
	FromEmail string `fig:"from_email,required"`
	FromName  string `fig:"from_name"`
}

type MailjetConfiger interface {
	MailjetConfig() MailjetConfig
}

type mailjetConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

func NewMailjetConfiger(getter kv.Getter) MailjetConfiger {
	return &mailjetConfiger{
		getter: getter,
	}
}

func (c *mailjetConfiger) MailjetConfig() MailjetConfig {
	return c.once.Do(func() interface{} {
		var cfg = MailjetConfig{
			FromName: defaultFromName,
		}

		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, configKeyMailjet)).
			With(figure.BaseHooks).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out mailjet config: %w", err))
		}

		return cfg
	}).(MailjetConfig)
}
