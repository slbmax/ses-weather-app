package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/slbmax/ses-weather-app/internal/api/ctx"
	"github.com/slbmax/ses-weather-app/internal/api/requests"
	"github.com/slbmax/ses-weather-app/internal/db"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewSubscribeRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		logger        = ctx.GetLogger(r)
		subscriptions = ctx.GetSubscriptions(r)
	)

	txErr := subscriptions.Transaction(func() error {
		token, _ := uuid.GenerateUUID()
		sub := db.Subscription{
			Email:     request.Email,
			City:      request.City,
			Frequency: request.Frequency,
			Token:     token,
			CreatedAt: time.Now(),
		}

		// writing data ahead to rollback in case of email sending failure
		// also, reducing the number of queries to the database
		if _, err := subscriptions.Insert(sub); err != nil {
			return err
		}

		// TODO: send confirmation email

		return nil
	})

	switch {
	case txErr == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(txErr, db.ErrSubscriptionExists):
		w.WriteHeader(http.StatusConflict)
	default:
		logger.WithError(txErr).Error("failed to insert subscription")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
