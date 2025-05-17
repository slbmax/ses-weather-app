package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/slbmax/ses-weather-app/assets/static"
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

var useMocks bool

func init() {
	runCmd.Flags().BoolVar(&useMocks, "mocks", false, "Use mock APIs for testing purposes")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Weather App server",
	RunE: func(cmd *cobra.Command, args []string) error {
		getter, err := kv.FromEnv()
		if err != nil {
			return fmt.Errorf("failed to get configer key-value getter: %w", err)
		}

		cfg := config.New(getter)

		stopCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		defer cancel()

		var (
			eg, ctx    = errgroup.WithContext(stopCtx)
			mail       mailer.Mailer
			weatherApi weatherapi.WeatherProvider
			logger     = cfg.Log()
		)

		if useMocks {
			mail = mailer.NewMockMailer()
			weatherApi = weatherapi.NewMockWeatherProvider()
		} else {
			mailjetCfg := cfg.MailjetConfig()
			mailjetClient := mailjet.NewClient(
				mailjetCfg.ApiKey,
				mailjetCfg.SecretKey,
				mailjet.From{
					Name:  mailjetCfg.FromName,
					Email: mailjetCfg.FromEmail,
				},
			)
			mail = mailer.NewMailer(mailjetClient)
			weatherApi = weatherapi.NewClient(cfg.WeatherAPIConfig().APIKey)
		}

		eg.Go(func() error {
			server := api.NewServer(
				cfg.Listener(),
				weatherApi,
				pg.NewDatabase(cfg.DB()),
				mail,
				logger.WithField("component", "api"),
			)

			return server.Run(ctx)
		})

		eg.Go(func() error {
			notificator.New(
				pg.NewDatabase(cfg.DB()),
				weatherApi,
				mail,
				logger.WithField("component", "notificator"),
			).Run(ctx)

			return nil
		})

		serveStaticCfg := cfg.ServeStaticConfig()
		if serveStaticCfg.Enabled {
			eg.Go(func() error {
				return static.Serve(ctx,
					static.IndexData{
						BaseApiUrl: serveStaticCfg.BaseApiUrl,
					},
					serveStaticCfg.Listener,
				)
			})
		}

		return eg.Wait()
	},
}
