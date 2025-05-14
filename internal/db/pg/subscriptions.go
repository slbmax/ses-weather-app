package pg

import (
	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/slbmax/ses-weather-app/internal/db"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	subscriptionsTable = "subscriptions"

	constraintUniqueEmail = "unique_email"
)

type subscriptionsQ struct {
	db *pgdb.DB
}

func NewSubscriptionsQ(db *pgdb.DB) db.Subscriptions {
	return &subscriptionsQ{
		db: db,
	}
}

func (s *subscriptionsQ) New() db.Subscriptions {
	return NewSubscriptionsQ(s.db.Clone())
}

func (s *subscriptionsQ) Insert(subscription db.Subscription) (id int64, err error) {
	stmt := squirrel.
		Insert(subscriptionsTable).
		SetMap(structs.Map(subscription)).
		Suffix("RETURNING id")

	err = s.db.Get(&id, stmt)
	if pgdb.IsConstraintErr(err, constraintUniqueEmail) {
		return 0, db.ErrSubscriptionExists
	}

	return
}

func (s *subscriptionsQ) Transaction(f func() error) error {
	return s.db.Transaction(f)
}
