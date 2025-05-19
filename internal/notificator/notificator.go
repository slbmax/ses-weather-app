package notificator

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/slbmax/ses-weather-app/internal/database"
	"github.com/slbmax/ses-weather-app/internal/mailer"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	"gitlab.com/distributed_lab/logan/v3"
)

const (
	notificatorInterval     = 30 * time.Second
	notificationParallelism = 10
)

type Notificator struct {
	db         database.Database
	weatherApi weatherapi.WeatherProvider
	mailer     mailer.Mailer
	logger     *logan.Entry
}

func New(db database.Database, weatherApi weatherapi.WeatherProvider, mailer mailer.Mailer, logger *logan.Entry) *Notificator {
	return &Notificator{
		db:         db,
		logger:     logger,
		mailer:     mailer,
		weatherApi: weatherApi,
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

// processPendingNotifications processes notifications in parallel
// it can (in a production env, must) be enhanced by using batch notification sending,
// bulk weather querying, and bulk updating, but, for this small project, it will be kept simple.
// Semaphore is used to limit the number of concurrent goroutines and possible rate limiting from third-party APIs
func (n *Notificator) processPendingNotifications(subs []database.Subscription) (processed int) {
	cache := newWeatherCache()
	semaphore := make(chan struct{}, notificationParallelism)
	successNotifications := new(atomic.Int32)

	wg := new(sync.WaitGroup)
	wg.Add(len(subs))
	for _, sub := range subs {
		semaphore <- struct{}{}
		go func(sub database.Subscription) {
			defer func() { <-semaphore; wg.Done() }()

			weather, ok := cache.Get(sub.City)
			if !ok {
				response, err := n.weatherApi.GetCurrentWeather(sub.City)
				if err != nil {
					n.logger.WithError(err).Errorf("failed to get weather for city %s", sub.City)
					return
				}
				weather = response.CurrentWeather
				cache.Set(sub.City, weather)
			}

			db := n.db.New()
			txErr := db.Transaction(func() error {
				if err := db.SubscriptionsQ().UpdateLastNotified(sub.Id, time.Now()); err != nil {
					return fmt.Errorf("failed to update last notified for id %v: %w", sub.Id, err)
				}

				if err := n.mailer.SendNotificationEmail(sub.Email, mailer.NotificationEmail{
					City:        sub.City,
					Temperature: weather.Temperature,
					Description: weather.Condition.Text,
					Humidity:    weather.Humidity,
					Frequency:   string(sub.Frequency),
				}); err != nil {
					return fmt.Errorf("failed to send notification email: %w", err)
				}

				return nil
			})
			if txErr != nil {
				n.logger.WithError(txErr).Error("failed to process notification")
				return
			}

			successNotifications.Add(1)
		}(sub)
	}
	wg.Wait()

	return int(successNotifications.Load())
}

type weatherCache struct {
	cache map[string]weatherapi.CurrentWeather
	mu    *sync.Mutex
}

func newWeatherCache() *weatherCache {
	return &weatherCache{
		cache: make(map[string]weatherapi.CurrentWeather),
		mu:    new(sync.Mutex),
	}
}

func (w *weatherCache) Get(city string) (weatherapi.CurrentWeather, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	weather, ok := w.cache[city]
	return weather, ok
}

func (w *weatherCache) Set(city string, weather weatherapi.CurrentWeather) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.cache[city] = weather
}
