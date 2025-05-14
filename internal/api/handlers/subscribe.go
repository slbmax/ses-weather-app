package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/slbmax/ses-weather-app/internal/api/ctx"
	"github.com/slbmax/ses-weather-app/internal/api/requests"
	"github.com/slbmax/ses-weather-app/internal/database"
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

	txErr := db.Transaction(func() error {
		// writing data ahead to rollback in case of email sending failure
		// also, reducing the number of queries to the database
		sub := database.Subscription{
			Email:     request.Email,
			City:      request.City,
			Frequency: request.Frequency,
			CreatedAt: time.Now(),
		}
		if sub.Id, err = db.SubscriptionsQ().Insert(sub); err != nil {
			return err
		}

		confirmationToken := database.Token{
			IsConfirmation: true,
			Token:          GenerateToken(),
			SubscriptionId: sub.Id,
		}
		unsubToken := database.Token{
			IsConfirmation: false,
			Token:          GenerateToken(),
			SubscriptionId: sub.Id,
		}
		if err = db.TokensQ().Insert(confirmationToken, unsubToken); err != nil {
			return err
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
		logger.WithError(txErr).Error("failed to insert subscription")
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
