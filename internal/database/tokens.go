package database

type TokensQ interface {
	New() TokensQ
	Insert(tokens ...Token) error
	Get(value string) (*Token, error)
	Delete(value string) error
}

type Token struct {
	Token          string `structs:"token" db:"token"`
	SubscriptionId int64  `structs:"subscription_id" db:"subscription_id"`
	IsConfirmation bool   `structs:"is_confirmation" db:"is_confirmation"`
}
