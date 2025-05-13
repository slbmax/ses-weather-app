package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/slbmax/ses-weather-app/internal/api/ctx"
	"github.com/slbmax/ses-weather-app/internal/api/handlers"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
)

type Server struct {
	logger     *logan.Entry
	listener   net.Listener
	weatherApi *weatherapi.Client
}

func NewServer(
	listener net.Listener,
	weatherApi *weatherapi.Client,
	logger *logan.Entry,
) *Server {
	return &Server{
		logger:     logger,
		listener:   listener,
		weatherApi: weatherApi,
	}
}

func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{Handler: s.requestHandler()}

	// graceful shutdown
	go func() {
		<-ctx.Done()

		shutdownDeadline, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownDeadline); err != nil {
			s.logger.WithError(err).Error("failed to shutdown http server")
		} else {
			s.logger.Info("http serving stopped: context canceled")
		}
	}()

	s.logger.Infof("http server listening on %s", s.listener.Addr().String())
	if err := srv.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) requestHandler() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.logger),
		ape.LoganMiddleware(s.logger),
		ape.CtxMiddleware(
			ctx.LoggerProvider(s.logger),
			ctx.WeatherApiProvider(s.weatherApi),
		),
	)

	r.Get("/weather", handlers.GetWeather)

	return r
}
