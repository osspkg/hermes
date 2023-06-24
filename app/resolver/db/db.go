package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/osspkg/go-sdk/orm"
	"github.com/osspkg/hermes-addons/dependency"
)

type DB struct {
	stmt orm.Stmt
}

func New(s orm.Stmt) dependency.ORM {
	return &DB{stmt: s}
}

func (v *DB) ExecContext(ctx context.Context, call func(q dependency.Executor)) error {
	return v.stmt.CallContext("", ctx, func(ctx context.Context, db *sql.DB) error {
		return callExecContext(ctx, db, call)
	})
}

func (v *DB) QueryContext(ctx context.Context, call func(q dependency.Querier)) error {
	return v.stmt.CallContext("", ctx, func(ctx context.Context, db *sql.DB) error {
		return callQueryContext(ctx, db, call)
	})
}

func (v *DB) TransactionContext(ctx context.Context, call func(v dependency.Tx)) error {
	q := &tx{}
	call(q)
	return v.stmt.TxContext("", ctx, func(ctx context.Context, tx *sql.Tx) error {
		for i, c := range q.v {
			if cc, ok := c.(func(q dependency.Executor)); ok {
				if err := callExecContext(ctx, tx, cc); err != nil {
					return err
				}
				continue
			}
			if cc, ok := c.(func(q dependency.Querier)); ok {
				if err := callQueryContext(ctx, tx, cc); err != nil {
					return err
				}
				continue
			}
			return fmt.Errorf("unknown query type #%d", i)
		}
		return nil
	})
}
