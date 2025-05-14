package pg

import (
	"github.com/Masterminds/squirrel"
	"github.com/slbmax/ses-weather-app/internal/database"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	tokensTable = "tokens"

	columnToken          = "token"
	columnSubscriptionId = "subscription_id"
	columnIsConfirmation = "is_confirmation"
)

type tokensQ struct {
	db *pgdb.DB
}

func NewTokensQ(db *pgdb.DB) database.TokensQ {
	return &tokensQ{db}
}

func (t *tokensQ) Insert(tokens ...database.Token) error {
	if len(tokens) == 0 {
		return nil
	}

	stmt := squirrel.Insert(tokensTable).
		Columns(columnToken, columnSubscriptionId, columnIsConfirmation)
	for _, token := range tokens {
		stmt = stmt.Values(token.Token, token.SubscriptionId, token.IsConfirmation)
	}

	return t.db.Exec(stmt)
}

func (t *tokensQ) New() database.TokensQ {
	return NewTokensQ(t.db.Clone())
}

func (t *tokensQ) Get(value string) (*database.Token, error) {
	return nil, nil
}

func (t *tokensQ) Delete(value string) error {
	return nil
}
