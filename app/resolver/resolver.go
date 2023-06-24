/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a LGPL-3.0 license that can be found in the LICENSE file.
 */

package resolver

import (
	"fmt"
	"sync"

	"github.com/osspkg/goppy/plugins/database"
	"github.com/osspkg/hermes-addons/dependency"
	db2 "github.com/osspkg/hermes/app/resolver/db"
)

type Resolver struct {
	data map[string]interface{}
	mux  sync.RWMutex
}

func New(db database.MySQL) *Resolver {
	obj := &Resolver{
		data: make(map[string]interface{}, 100),
	}
	obj.data[dependency.Database] = db2.New(db.Pool("main"))
	return obj
}

func (v *Resolver) Get(dep string) (interface{}, error) {
	v.mux.RLock()
	defer v.mux.RUnlock()
	if obj, ok := v.data[dep]; ok {
		return obj, nil
	}
	return nil, fmt.Errorf("dependency not found: %s", dep)
}

func (v *Resolver) Set(dep string, obj interface{}) error {
	v.mux.Lock()
	defer v.mux.Unlock()
	if _, ok := v.data[dep]; ok {
		return fmt.Errorf("dependency already exist: %s", dep)
	}
	return nil
}
