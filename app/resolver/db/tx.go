package db

import "github.com/osspkg/hermes-addons/dependency"

type tx struct {
	v []interface{}
}

func (v *tx) Exec(vv ...func(q dependency.Executor)) {
	for _, f := range vv {
		v.v = append(v.v, f)
	}
}

func (v *tx) Query(vv ...func(q dependency.Querier)) {
	for _, f := range vv {
		v.v = append(v.v, f)
	}
}
