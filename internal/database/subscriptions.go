package database

import (
	"errors"
	"time"
)

type SubscriptionFrequency string

func (f SubscriptionFrequency) Valid() bool {
	switch f {
	case SubscriptionFrequencyDaily, SubscriptionFrequencyHourly:
		return true
	default:
		return false
	}
}

const (
	SubscriptionFrequencyDaily  SubscriptionFrequency = "daily"
	SubscriptionFrequencyHourly SubscriptionFrequency = "hourly"
)

var (
	ErrSubscriptionExists = errors.New("subscription already exists")
)

type SubscriptionsQ interface {
	// New creates a new instance of SubscriptionsQ (separate conn)
	New() SubscriptionsQ
	Insert(subscription Subscription) (id int64, err error)
	//MarkConfirmed(id int64) error
	//Delete(id int64) error
	//UpdateLastNotified(id int64, lastNotifiedAt time.Time) error
}

type Subscription struct {
	Id             int64                 `structs:"-" db:"id"`
	Email          string                `structs:"email" db:"email"`
	City           string                `structs:"city" db:"city"`
	Frequency      SubscriptionFrequency `structs:"frequency" db:"frequency"`
	Confirmed      bool                  `structs:"confirmed" db:"confirmed"`
	CreatedAt      time.Time             `structs:"created_at" db:"created_at"`
	LastNotifiedAt *time.Time            `structs:"last_notified_at" db:"last_notified_at"`
}
