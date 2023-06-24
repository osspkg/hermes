package db

import (
	"context"

	"github.com/osspkg/hermes-addons/dependency"
)

type exec struct {
	Q string
	P [][]interface{}
	B func(rowsAffected, lastInsertId int64) error
}

func (v *exec) SQL(query string, args ...interface{}) {
	v.Q = query
	v.Params(args...)
}

func (v *exec) Params(args ...interface{}) {
	if len(args) > 0 {
		v.P = append(v.P, args)
	}
}
func (v *exec) Bind(call func(rowsAffected, lastInsertId int64) error) {
	v.B = call
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func callExecContext(ctx context.Context, db dbi, call func(q dependency.Executor)) error {
	q := &exec{}
	call(q)
	if len(q.P) == 0 {
		q.P = append(q.P, []interface{}{})
	}
	stmt, err := db.PrepareContext(ctx, q.Q)
	if err != nil {
		return err
	}
	defer stmt.Close() //nolint: errcheck
	var rowsAffected, lastInsertId int64
	for _, params := range q.P {
		result, err0 := stmt.ExecContext(ctx, params...)
		if err0 != nil {
			return err0
		}
		rows, err0 := result.RowsAffected()
		if err0 != nil {
			return err0
		}
		rowsAffected += rows
		rows, err0 = result.LastInsertId()
		if err0 != nil {
			return err0
		}
		lastInsertId = rows
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	if q.B == nil {
		return nil
	}
	return q.B(rowsAffected, lastInsertId)
}
