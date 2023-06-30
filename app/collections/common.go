package collections

import (
	"github.com/osspkg/go-sdk/orm"
	"github.com/osspkg/goppy/plugins/database"
)

type Collections struct {
	stmt orm.Stmt
}

func New(db database.MySQL) *Collections {
	return &Collections{
		stmt: db.Pool("main"),
	}
}
