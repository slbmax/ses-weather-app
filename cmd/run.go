package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/slbmax/ses-weather-app/internal/api"
	"github.com/slbmax/ses-weather-app/internal/config"
	"github.com/slbmax/ses-weather-app/internal/database/pg"
	"github.com/slbmax/ses-weather-app/internal/mailer"
	"github.com/slbmax/ses-weather-app/internal/notificator"
	"github.com/slbmax/ses-weather-app/pkg/mailjet"
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
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		defer cancel()
		eg, ctx := errgroup.WithContext(ctx)

		mailjetCfg := cfg.MailjetConfig()
		mailjetClient := mailjet.NewClient(
			mailjetCfg.ApiKey,
			mailjetCfg.SecretKey,
			mailjet.From{
				Name:  mailjetCfg.FromName,
				Email: mailjetCfg.FromEmail,
			},
		)

		mailer := mailer.NewMailer(mailjetClient)
		weatherApi := weatherapi.NewClient(cfg.WeatherAPIConfig().APIKey)
		logger := cfg.Log()

		eg.Go(func() error {
			server := api.NewServer(
				cfg.Listener(),
				weatherApi,
				pg.NewDatabase(cfg.DB()),
				mailer,
				logger.WithField("component", "api"),
			)

			return server.Run(ctx)
		})

		eg.Go(func() error {
			notificator.New(
				pg.NewDatabase(cfg.DB()),
				weatherApi,
				mailer,
				logger.WithField("component", "notificator"),
			).Run(ctx)

			return nil
		})

		return eg.Wait()
	},
}
