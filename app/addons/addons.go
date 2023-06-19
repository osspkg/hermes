package addons

import (
	"fmt"
	"github.com/osspkg/go-algorithms/graph/kahn"
	"github.com/osspkg/go-sdk/app"
	"github.com/osspkg/go-sdk/errors"
	"github.com/osspkg/go-sdk/iofile"
	"github.com/osspkg/go-sdk/log"
	"github.com/osspkg/hermes/app/pkg/util"
	"plugin"
	"sync"

	"github.com/osspkg/hermes-addons/api1"
)

const (
	Symbol     = "HermesAPI"
	FileNameSO = "addon.so"
)

type Addons struct {
	addons       map[string]api1.Api
	system       map[string]interface{}
	dependencies []string

	conf *Config
	mux  sync.RWMutex
}

func New(c *Config) *Addons {
	return &Addons{
		addons:       make(map[string]api1.Api, 100),
		system:       make(map[string]interface{}, 100),
		dependencies: make([]string, 0, 100),
		conf:         c,
	}
}

func (v *Addons) Down() error {
	v.mux.Lock()
	defer v.mux.Unlock()

	var errs error
	for name, api := range v.addons {
		if err := api.Down(); err != nil {
			errs = errors.Wrapf(errs, "%s: %w", name, err)
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

	if err = v.resolve(); err != nil {
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

func (v *Addons) resolve() error {
	graph := kahn.New()

	v.mux.Lock()
	defer v.mux.Unlock()

	for cur, api := range v.addons {
		if err := v.initSystemDependency(api, api.Dependency()); err != nil {
			return fmt.Errorf("init addon system dependency for `%s`: %w", cur, err)
		}

		for _, dep := range api.Dependency() {
			if err := graph.Add(cur, dep); err != nil {
				return fmt.Errorf("add addon dependency to graph `%s=>%s`: %w", cur, dep, err)
			}
		}
	}

	if err := graph.Build(); err != nil {
		return fmt.Errorf("build addon dependency graph: %w", err)
	}

	v.dependencies = graph.Result()

	for _, dep := range v.dependencies {
		if _, ok := v.addons[dep]; !ok {
			return fmt.Errorf("addon dependency not found: %s", dep)
		}
	}

	util.FlipStringSlice(v.dependencies)

	return nil
}

func (v *Addons) initSystemDependency(api api1.Api, deps []string) error {
	for _, dep := range deps {
		switch dep {

		default:
			fmt.Println("->", dep)
		}
	}
	return nil
}
