package database

type Database interface {
	New() Database
	SubscriptionsQ() SubscriptionsQ
	Transaction(func() error) error
}
