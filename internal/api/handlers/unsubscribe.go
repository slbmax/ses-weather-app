package handlers

import (
	"errors"
	"net/http"

	"github.com/slbmax/ses-weather-app/internal/api/ctx"
	"github.com/slbmax/ses-weather-app/internal/api/requests"
	"github.com/slbmax/ses-weather-app/internal/database"
)

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUnsubscribeRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		log = ctx.GetLogger(r)
		db  = ctx.GetDatabase(r)
	)

	// as the specification states, we can delete the subscription by using
	// - confirmation token at the unconfirmed state
	// - unsubscribe token at the confirmed state
	// (Unsubscribes an email from weather updates using the token sent in email__S__)
	err = db.SubscriptionsQ().DeleteByToken(request.Token)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, database.ErrNoRowsAffected):
		w.WriteHeader(http.StatusNotFound)
	default:
		log.WithError(err).Error("failed delete subscription")
		w.WriteHeader(http.StatusInternalServerError)
	}

}
