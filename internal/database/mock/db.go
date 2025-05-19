package mock

import "github.com/slbmax/ses-weather-app/internal/database"

type db struct {
	subscriptionsMock *MockSubscriptionsQ
}

func NewDatabase(mock *MockSubscriptionsQ) database.Database {
	return &db{
		subscriptionsMock: mock,
	}
}

func (d *db) New() database.Database {
	return &db{
		subscriptionsMock: d.subscriptionsMock,
	}
}

func (d *db) SubscriptionsQ() database.SubscriptionsQ {
	return d.subscriptionsMock
}

func (d *db) Transaction(fn func() error) error {
	return fn()
}
