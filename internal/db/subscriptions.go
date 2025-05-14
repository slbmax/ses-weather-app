package db

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

type Subscriptions interface {
	// New creates a new instance of Subscriptions (separate conn)
	New() Subscriptions
	Insert(subscription Subscription) (id int64, err error)
	//MarkConfirmed(id int64) error

	Transaction(f func() error) error
}

type Subscription struct {
	Id             int64                 `structs:"-" db:"id"`
	Email          string                `structs:"email" db:"email"`
	City           string                `structs:"city" db:"city"`
	Frequency      SubscriptionFrequency `structs:"frequency" db:"frequency"`
	Token          string                `structs:"token" db:"token"`
	Confirmed      bool                  `structs:"confirmed" db:"confirmed"`
	CreatedAt      time.Time             `structs:"created_at" db:"created_at"`
	LastNotifiedAt *time.Time            `structs:"last_notified_at" db:"last_notified_at"`
}
