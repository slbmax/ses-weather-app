package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/slbmax/ses-weather-app/internal/api/ctx"
	"github.com/slbmax/ses-weather-app/internal/api/requests"
	"github.com/slbmax/ses-weather-app/internal/database"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
)

const (
	tokenLengthRaw = 16
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewSubscribeRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		logger = ctx.GetLogger(r)
		db     = ctx.GetDatabase(r)
	)

	_, err = ctx.GetWeatherClient(r).GetCurrentWeather(request.City)
	if err != nil {
		if errors.Is(err, weatherapi.ErrCityNotFound) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			logger.WithError(err).Error("failed to get current weather")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	txErr := db.Transaction(func() error {
		// writing data ahead to rollback in case of email sending failure
		// also, reducing the number of queries to the database
		sub := database.Subscription{
			Email:     request.Email,
			City:      request.City,
			Token:     GenerateToken(),
			Frequency: request.Frequency,
			CreatedAt: time.Now(),
		}
		if sub.Id, err = db.SubscriptionsQ().Insert(sub); err != nil {
			return fmt.Errorf("failed to insert subscription: %w", err)
		}

		// TODO: send confirmation email

		return nil
	})

	switch {
	case txErr == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(txErr, database.ErrSubscriptionExists):
		w.WriteHeader(http.StatusConflict)
	default:
		logger.WithError(txErr).Error("failed to execute transaction")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GenerateToken() string {
	b := make([]byte, tokenLengthRaw)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate token: " + err.Error())
	}

	return hex.EncodeToString(b)
}
