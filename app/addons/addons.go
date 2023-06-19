/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a LGPL-3.0 license that can be found in the LICENSE file.
 */

package addons

import (
	"fmt"
	"plugin"
	"sync"

	"github.com/osspkg/go-algorithms/graph/kahn"
	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/errors"
	"github.com/osspkg/go-sdk/iofile"
	"github.com/osspkg/go-sdk/log"
	"github.com/osspkg/hermes-addons/api1"
	"github.com/osspkg/hermes/app/pkg/util"
	"github.com/osspkg/hermes/app/resolver"
)

const (
	Symbol     = "HermesAPI"
	FileNameSO = "addon.so"
)

type Addons struct {
	addons       map[string]api1.Api
	dependencies []string

	resolver *resolver.Resolver
	conf     *Config
	mux      sync.RWMutex
}

func New(c *Config, r *resolver.Resolver) *Addons {
	return &Addons{
		addons:       make(map[string]api1.Api, 100),
		dependencies: make([]string, 0, 100),
		resolver:     r,
		conf:         c,
	}
}

func (v *Addons) Down() error {
	v.mux.Lock()
	defer v.mux.Unlock()

	var errs error
	for _, dep := range v.dependencies {
		api, ok := v.addons[dep]
		if !ok {
			continue
		}
		if err := api.Down(); err != nil {
			errs = errors.Wrapf(errs, "addon stop `%s`: %w", dep, err)
		}
	}
	return errs
}

func (v *Addons) Up(ctx app.Context) error {
	files, err := iofile.Search(v.conf.Addons, FileNameSO)
	if err != nil {
		return fmt.Errorf("load addons from `%s`: %w", v.conf.Addons, err)
	}
	for _, filename := range files {
		if err = v.load(filename); err != nil {
			return err
		}
	}

	if err = v.resolve(ctx); err != nil {
		return err
	}

	return nil
}

func (v *Addons) load(filename string) error {
	mod, err := plugin.Open(filename)
	if err != nil {
		return err
	}

	symApi, err := mod.Lookup(Symbol)
	if err != nil {
		return err
	}

	apiInit, ok := symApi.(func() api1.Api)
	if !ok {
		return fmt.Errorf("invalid api v1 for `%s`", filename)
	}

	api := apiInit()

	v.mux.Lock()
	defer v.mux.Unlock()

	log.WithField(api.PkgName(), filename).Infof("Load addon")

	v.addons[api.PkgName()] = api

	return nil
}

func (v *Addons) resolve(ctx app.Context) error {
	graph := kahn.New()

	v.mux.Lock()
	defer v.mux.Unlock()

	for cur, api := range v.addons {
		for _, dep := range api.Dependency() {
			if err := graph.Add(cur, dep); err != nil {
				return fmt.Errorf("add addon dependency to graph `%s=>%s`: %w", cur, dep, err)
			}
		}
	}

	if err := graph.Build(); err != nil {
		return fmt.Errorf("build addon dependency graph: %w", err)
	}

	deps := graph.Result()

	for _, dep := range deps {
		if v.resolver.Has(dep) {
			continue
		}
		if api, ok := v.addons[dep]; ok {
			if err := api.Inject(v.resolver); err != nil {
				return fmt.Errorf("addon init `%s`: %w", dep, err)
			}
			if err := v.resolver.Set(dep, api); err != nil {
				return fmt.Errorf("addon save to resolver `%s`: %w", dep, err)
			}
			v.dependencies = append(v.dependencies, dep)
			if err := api.Up(ctx.Context()); err != nil {
				return fmt.Errorf("addon start `%s`: %w", dep, err)
			}
			continue
		}
		return fmt.Errorf("addon dependency not found: %s", dep)
	}

	util.FlipStringSlice(v.dependencies)

	return nil
}
