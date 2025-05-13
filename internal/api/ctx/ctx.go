package ctx

import (
	"context"
	"net/http"

	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	ctxKeyLogger ctxKey = iota
	ctxKeyWeatherApi
)

func LoggerProvider(l *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKeyLogger, l)
	}
}
func GetLogger(r *http.Request) *logan.Entry {
	return r.Context().Value(ctxKeyLogger).(*logan.Entry)
}

func WeatherApiProvider(api *weatherapi.Client) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKeyWeatherApi, api)
	}
}

func GetWeatherClient(r *http.Request) *weatherapi.Client {
	return r.Context().Value(ctxKeyWeatherApi).(*weatherapi.Client)
}
