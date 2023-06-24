package db

import (
	"context"

	"github.com/osspkg/hermes-addons/dependency"
)

type query struct {
	Q string
	P []interface{}
	B func(bind dependency.Scanner) error
}

func (v *query) SQL(query string, args ...interface{}) {
	v.Q, v.P = query, args
}

func (v *query) Bind(call func(bind dependency.Scanner) error) {
	v.B = call
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func callQueryContext(ctx context.Context, db dbi, call func(q dependency.Querier)) error {
	q := &query{}
	call(q)
	rows, err := db.QueryContext(ctx, q.Q, q.P...)
	if err != nil {
		return err
	}
	defer rows.Close() //nolint: errcheck
	if q.B != nil {
		for rows.Next() {
			if err = q.B(rows); err != nil {
				return err
			}
		}
	}
	if err = rows.Close(); err != nil {
		return err
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
