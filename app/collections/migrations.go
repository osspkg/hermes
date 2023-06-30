package collections

import (
	"context"
	"fmt"

	"github.com/osspkg/go-sdk/orm"
)

func (v *Collections) ApplyMigrations(ctx context.Context, addon, name, query string) error {
	var has bool

	err := v.stmt.QueryContext("", ctx, func(q orm.Querier) {
		q.SQL("SELECT id FROM `hermes_addon_migrations` WHERE `addon` = ? AND `data` = ? LIMIT 1;", addon, name)
		q.Bind(func(bind orm.Scanner) error {
			var id int
			if err := bind.Scan(&id); err != nil {
				return err
			}
			if id > 0 {
				has = true
			}
			return nil
		})
	})
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	err = v.stmt.ExecContext("", ctx, func(q orm.Executor) {
		q.SQL(query)
	})
	if err != nil {
		return err
	}

	return v.stmt.ExecContext("", ctx, func(q orm.Executor) {
		q.SQL("INSERT INTO `hermes_addon_migrations` (`addon`, `data`, `created_at`) VALUES (?, ?, now());")
		q.Params(addon, name)
		q.Bind(func(rowsAffected, lastInsertId int64) error {
			if lastInsertId == 0 {
				return fmt.Errorf("fail")
			}
			return nil
		})
	})
}
