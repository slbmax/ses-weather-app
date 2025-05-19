package pg

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/slbmax/ses-weather-app/internal/database"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	subscriptionsTable = "subscriptions"

	columnId             = "id"
	columnConfirmed      = "confirmed"
	columnToken          = "token"
	columnLastNotifiedAt = "last_notified_at"

	constraintUniqueEmail = "unique_email"
)

type subscriptionsQ struct {
	db *pgdb.DB
}

func NewSubscriptionsQ(db *pgdb.DB) database.SubscriptionsQ {
	return &subscriptionsQ{
		db: db,
	}
}

func (s *subscriptionsQ) New() database.SubscriptionsQ {
	return NewSubscriptionsQ(s.db.Clone())
}

func (s *subscriptionsQ) Insert(subscription database.Subscription) (id int64, err error) {
	stmt := squirrel.
		Insert(subscriptionsTable).
		SetMap(structs.Map(subscription)).
		Suffix("RETURNING id")

	err = s.db.Get(&id, stmt)
	if pgdb.IsConstraintErr(err, constraintUniqueEmail) {
		return 0, database.ErrSubscriptionExists
	}

	return
}

func (s *subscriptionsQ) GetByToken(token string) (*database.Subscription, error) {
	stmt := squirrel.
		Select("*").
		From(subscriptionsTable).
		Where(squirrel.Eq{columnToken: token})

	var subscription database.Subscription
	err := s.db.Get(&subscription, stmt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &subscription, err
}

func (s *subscriptionsQ) UpdateConfirmed(id int64, unsubscribeToken string) error {
	stmt := squirrel.
		Update(subscriptionsTable).
		Set(columnConfirmed, true).
		Set(columnToken, unsubscribeToken).
		Where(squirrel.Eq{
			columnId:        id,
			columnConfirmed: false, // generally can be omitted, but still
		}).Suffix("RETURNING *")

	// actually, the "WITH rows as (UPDATE ... RETURNING *) SELECT COUNT(*) FROM rows" allows us
	// to bypass the forbidden aggregated functions in the "RETURNING" clause, but for simplicity
	// we just use the "RETURNING *" clause and check the length of the result

	var affectedRows []database.Subscription
	if err := s.db.Select(&affectedRows, stmt); err != nil {
		return err
	} else if len(affectedRows) == 0 {
		return database.ErrNoRowsAffected
	}

	return nil
}

func (s *subscriptionsQ) DeleteByToken(token string) error {
	stmt := squirrel.
		Delete(subscriptionsTable).
		Where(squirrel.Eq{columnToken: token})

	if result, err := s.db.ExecWithResult(stmt); err != nil {
		return err
	} else if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return database.ErrNoRowsAffected
	} else {
		return nil
	}
}

func (s *subscriptionsQ) SelectToNotify() ([]database.Subscription, error) {
	stmt := squirrel.
		Select("*").
		From(subscriptionsTable).
		Where(squirrel.Eq{columnConfirmed: true}).
		Where(squirrel.Or{
			squirrel.Eq{columnLastNotifiedAt: nil},
			squirrel.Expr(`
                last_notified_at <= CURRENT_TIMESTAMP - (
					CASE frequency
						WHEN 'hourly' THEN INTERVAL '1 hour'
						ELSE INTERVAL '24 hours'
					END
				)
            `),
		})

	var subscriptions []database.Subscription
	if err := s.db.Select(&subscriptions, stmt); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *subscriptionsQ) UpdateLastNotified(id int64, lastNotifiedAt time.Time) error {
	stmt := squirrel.
		Update(subscriptionsTable).
		Set(columnLastNotifiedAt, lastNotifiedAt).
		Where(squirrel.Eq{columnId: id})

	return s.db.Exec(stmt)
}
