/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a LGPL-3.0 license that can be found in the LICENSE file.
 */

package addons

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/iofile"
	"github.com/osspkg/go-sdk/log"
	hermesaddons "github.com/osspkg/hermes-addons"
	"github.com/osspkg/hermes/app/resolver"
)

const (
	Symbol   = "HermesAPI"
	Manifest = "manifest.json"
)

type (
	Addons struct {
		addons map[string]*Addon
		status *Status

		resolver *resolver.Resolver
		conf     *Config
		mux      sync.RWMutex
	}
	Addon struct {
		API      hermesaddons.Api
		Manifest ManifestModel
	}
)

func New(c *Config, r *resolver.Resolver) *Addons {
	return &Addons{
		addons:   make(map[string]*Addon, 100),
		status:   NewStatus(),
		resolver: r,
		conf:     c,
	}
}

func (v *Addons) Down() error {
	v.mux.Lock()
	defer v.mux.Unlock()

	for _, addon := range v.addons {
		if err := addon.API.Down(); err != nil {
			log.WithFields(log.Fields{
				"pkg": addon.Manifest.PkgName,
				"ver": addon.Manifest.Version,
				"err": err.Error(),
			}).Errorf("Unload addon")
		}
	}
	return nil
}

func (v *Addons) Up(ctx app.Context) error {
	return v.ReloadAll(ctx.Context())
}

func (v *Addons) Available() ([]ManifestModel, error) {
	files, err := iofile.Search(v.conf.Addons, Manifest)
	if err != nil {
		return nil, fmt.Errorf("load addons from `%s`: %w", v.conf.Addons, err)
	}
	result := make([]ManifestModel, 0, len(files))
	for _, filename := range files {
		b, err := os.ReadFile(filename)
		if err != nil {
			log.WithFields(log.Fields{
				"filename": filename,
				"err":      err.Error(),
			}).Errorf("Read " + Manifest)
			continue
		}
		model := ManifestModel{}
		if err = json.Unmarshal(b, &model); err != nil {
			log.WithFields(log.Fields{
				"filename": filename,
				"err":      err.Error(),
			}).Errorf("Unmarshal " + Manifest)
			continue
		}
		model.Filename = filepath.Dir(filename) + "/" + model.Filename
		result = append(result, model)
	}
	return result, nil
}

func (v *Addons) ReloadAll(ctx context.Context) error {
	models, err := v.Available()
	if err != nil {
		return err
	}
	for _, model := range models {
		if err = v.Load(ctx, model); err != nil {
			log.WithFields(log.Fields{
				"pkg":  model.PkgName,
				"ver":  model.Version,
				"file": model.Filename,
				"err":  err.Error(),
			}).Errorf("Load addon")
		}
	}
	return nil
}

func (v *Addons) Load(ctx context.Context, model ManifestModel) error {
	mod, err := plugin.Open(model.Filename)
	if err != nil {
		return err
	}

	symApi, err := mod.Lookup(Symbol)
	if err != nil {
		return err
	}

	apiInit, ok := symApi.(func() hermesaddons.Api)
	if !ok {
		return fmt.Errorf("invalid api v1")
	}

	addon := apiInit()

	if err = addon.Inject(v.resolver); err != nil {
		return fmt.Errorf("init addon: %w", err)
	}

	_ = v.Unload(model.PkgName)

	v.mux.Lock()
	defer v.mux.Unlock()

	v.addons[model.PkgName] = &Addon{
		API:      addon,
		Manifest: model,
	}

	log.WithFields(log.Fields{
		"pkg":  model.PkgName,
		"ver":  model.Version,
		"file": model.Filename,
	}).Infof("Load addon")

	err = addon.Up(ctx)

	v.status.Set(model, err)

	return err
}

func (v *Addons) Unload(pkgName string) error {
	v.mux.Lock()
	defer v.mux.Unlock()

	addon, ok := v.addons[pkgName]
	if !ok {
		return nil
	}
	err := addon.API.Down()
	if err != nil {
		log.WithFields(log.Fields{
			"pkg": addon.Manifest.PkgName,
			"ver": addon.Manifest.Version,
			"err": err.Error(),
		}).Errorf("Unload addon")
		return err
	}
	return nil
}
