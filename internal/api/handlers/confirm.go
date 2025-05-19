package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/slbmax/ses-weather-app/internal/api/ctx"
	"github.com/slbmax/ses-weather-app/internal/api/requests"
	"github.com/slbmax/ses-weather-app/internal/database"
	"github.com/slbmax/ses-weather-app/internal/mailer"
)

var (
	errSubscriptionConfirmed = errors.New("subscription already confirmed")
)

func Confirm(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewConfirmRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		logger = ctx.GetLogger(r)
		db     = ctx.GetDatabase(r)
		mail   = ctx.GetMailer(r)
	)

	txErr := db.Transaction(func() error {
		subscription, err := db.SubscriptionsQ().GetByToken(request.Token)
		if err != nil {
			return fmt.Errorf("failed to get subcription: %w", err)
		} else if subscription == nil {
			return sql.ErrNoRows
		} else if subscription.Confirmed {
			return errSubscriptionConfirmed
		}

		unsubToken := GenerateToken()
		if err = db.SubscriptionsQ().UpdateConfirmed(subscription.Id, unsubToken); err != nil {
			return fmt.Errorf("failed to confirm subscription: %w", err)
		}

		if err = mail.SendConfirmationSuccessEmail(subscription.Email, mailer.ConfirmationSuccessEmail{
			Token:     unsubToken,
			City:      subscription.City,
			Frequency: string(subscription.Frequency),
		}); err != nil {
			return fmt.Errorf("failed to send confirmation success email: %w", err)
		}

		// additionally, the notification email can be sent immediately,
		// for simplicity, leaving this to the notifier

		return nil
	})

	switch {
	case txErr == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, errSubscriptionConfirmed):
		// this token is supposed to be used as an unsubscribe token
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(txErr, database.ErrNoRowsAffected) ||
		errors.Is(txErr, sql.ErrNoRows):
		w.WriteHeader(http.StatusNotFound)
	default:
		logger.WithError(txErr).Error("failed to execute transaction")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
