package notificator

import (
	"context"
	"fmt"
	"time"

	"github.com/slbmax/ses-weather-app/internal/database"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	notificatorInterval = 30 * time.Second
)

type Notificator struct {
	db         database.Database
	weatherApi *weatherapi.Client
	logger     *logan.Entry
}

func New(db database.Database, logger *logan.Entry) *Notificator {
	return &Notificator{
		db:     db,
		logger: logger,
	}
}

func (n *Notificator) Run(ctx context.Context) {
	ticker := time.NewTicker(notificatorInterval)
	defer ticker.Stop()

	// trick to bypass initial tick delay
	for ; ; waitForTickerOrCtx(ctx, ticker) {
		if isCancelled(ctx) {
			n.logger.Info("notificator stopped")
			return
		}

		subs, err := n.db.SubscriptionsQ().SelectToNotify()
		if err != nil {
			n.logger.WithError(err).Error("failed to select subscriptions to notify")
			continue
		} else if len(subs) == 0 {
			n.logger.Info("no subscriptions to notify, sleeping...")
			continue
		}

		n.logger.Infof("got %v notifications to process", len(subs))
		processed := n.processPendingNotifications(subs)
		n.logger.Infof("successfully processed %v notifications", processed)
	}
}

func (n *Notificator) processPendingNotifications(subs []database.Subscription) (successNotifications int) {
	weatherCache := make(map[string]weatherapi.CurrentWeather)

	for _, sub := range subs {
		weather, ok := weatherCache[sub.City]
		if !ok {
			var err error
			response, err := n.weatherApi.GetCurrentWeather(sub.City)
			if err != nil {
				n.logger.WithError(err).Errorf("failed to get weather for city %s", sub.City)
				continue
			}

			weather = response.CurrentWeather
			weatherCache[sub.City] = weather
		}

		// it can (in a production env, must) be enhanced by using batch notification sending and bulk updating,
		// but, for this small project, it will be kept simple
		txErr := n.db.Transaction(func() error {
			if err := n.db.SubscriptionsQ().UpdateLastNotified(sub.Id, time.Now()); err != nil {
				return fmt.Errorf("failed to update last notified for id %v: %w", sub.Id, err)
			}

			// TODO: send notification

			return nil
		})
		if txErr != nil {
			n.logger.WithError(txErr).Error("failed to process notification")
			continue
		}

		successNotifications++
	}

	return
}
