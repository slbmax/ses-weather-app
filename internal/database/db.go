package database

type Database interface {
	New() Database
	SubscriptionsQ() SubscriptionsQ
	TokensQ() TokensQ
	Transaction(func() error) error
}
