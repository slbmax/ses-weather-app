package ctx

import (
	"context"
	"net/http"

	"github.com/slbmax/ses-weather-app/internal/db"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	ctxKeyLogger ctxKey = iota
	ctxKeyWeatherApi
	ctxKeySubscriptions
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

func SubscriptionsProvider(subscriptions db.Subscriptions) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKeySubscriptions, subscriptions)
	}
}

func GetSubscriptions(r *http.Request) db.Subscriptions {
	return r.Context().Value(ctxKeySubscriptions).(db.Subscriptions).New() // use New() to get a separate conn
}
