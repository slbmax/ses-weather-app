package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/slbmax/ses-weather-app/internal/api"
	"github.com/slbmax/ses-weather-app/internal/config"
	"github.com/slbmax/ses-weather-app/internal/db/pg"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/kit/kv"
	"golang.org/x/sync/errgroup"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Weather App server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.New(kv.MustFromEnv())

		eg := errgroup.Group{}
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		defer cancel()

		weatherApi := weatherapi.NewClient(cfg.WeatherAPIConfig().APIKey)
		logger := cfg.Log()
		server := api.NewServer(
			cfg.Listener(),
			weatherApi,
			pg.NewSubscriptionsQ(cfg.DB()),
			logger.WithField("component", "api"),
		)

		eg.Go(func() error {
			return server.Run(ctx)
		})

		return eg.Wait()
	},
}
