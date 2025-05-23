package ctx

import (
	"context"
	"net/http"

	"github.com/slbmax/ses-weather-app/internal/database"
	"github.com/slbmax/ses-weather-app/internal/mailer"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	ctxKeyLogger ctxKey = iota
	ctxKeyWeatherApi
	ctxKeyDatabase
	ctxKeyMailer
)

func LoggerProvider(l *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKeyLogger, l)
	}
}
func GetLogger(r *http.Request) *logan.Entry {
	return r.Context().Value(ctxKeyLogger).(*logan.Entry)
}

func WeatherApiProvider(api weatherapi.WeatherProvider) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKeyWeatherApi, api)
	}
}

func GetWeatherClient(r *http.Request) weatherapi.WeatherProvider {
	return r.Context().Value(ctxKeyWeatherApi).(weatherapi.WeatherProvider)
}

func DatabaseProvider(db database.Database) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKeyDatabase, db)
	}
}

func GetDatabase(r *http.Request) database.Database {
	return r.Context().Value(ctxKeyDatabase).(database.Database).New() // returns unique connection (for transaction purposes)
}

func MailerProvider(mailer mailer.Mailer) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxKeyMailer, mailer)
	}
}

func GetMailer(r *http.Request) mailer.Mailer {
	return r.Context().Value(ctxKeyMailer).(mailer.Mailer)
}
