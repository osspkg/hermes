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
)

type Resolver struct {
	system map[string]interface{}
	addons map[string]interface{}
	mux    sync.RWMutex
}

func New(
	db database.MySQL,
) *Resolver {
	obj := &Resolver{
		system: make(map[string]interface{}, 100),
		addons: make(map[string]interface{}, 100),
	}
	obj.system[dependency.Database] = db.Pool("main")
	return obj
}

func (v *Resolver) Has(dep string) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()
	if _, ok := v.system[dep]; ok {
		return true
	}
	if _, ok := v.addons[dep]; ok {
		return true
	}
	return false
}

func (v *Resolver) Get(dep string) (interface{}, error) {
	v.mux.RLock()
	defer v.mux.RUnlock()
	if obj, ok := v.system[dep]; !ok {
		return obj, nil
	}
	if obj, ok := v.addons[dep]; !ok {
		return obj, nil
	}
	return nil, fmt.Errorf("dependency not found: %s", dep)
}

func (v *Resolver) Set(dep string, obj interface{}) error {
	v.mux.Lock()
	defer v.mux.Unlock()
	if _, ok := v.system[dep]; ok {
		return fmt.Errorf("dependency already exist: %s", dep)
	}
	if _, ok := v.addons[dep]; ok {
		return fmt.Errorf("dependency already exist: %s", dep)
	}
	v.addons[dep] = obj
	return nil
}
