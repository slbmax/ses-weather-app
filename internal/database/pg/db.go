package pg

import (
	"github.com/slbmax/ses-weather-app/internal/database"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type db struct {
	db *pgdb.DB
}

func NewDatabase(database *pgdb.DB) database.Database {
	return &db{database}
}

func (d *db) New() database.Database {
	return NewDatabase(d.db.Clone())
}

func (d *db) SubscriptionsQ() database.SubscriptionsQ {
	return NewSubscriptionsQ(d.db)
}

func (d *db) TokensQ() database.TokensQ {
	return NewTokensQ(d.db)
}

func (d *db) Transaction(fn func() error) error {
	return d.db.Transaction(fn)
}
