package addons

import (
	"sync"
)

type (
	Status struct {
		data []StatusItem
		mux  sync.RWMutex
	}
	StatusItem struct {
		Manifest ManifestModel
		Err      error
	}
)

func NewStatus() *Status {
	return &Status{
		data: make([]StatusItem, 0, 10),
	}
}

func (v *Status) Set(model ManifestModel, err error) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.data = append(v.data, StatusItem{
		Manifest: model,
		Err:      err,
	})
}

func (v *Status) Reset() {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.data = v.data[:0]
}

func (v *Status) List() []StatusItem {
	v.mux.RLock()
	defer v.mux.RUnlock()

	pkgItems := make([]StatusItem, 0, len(v.data))
	pkgItems = append(pkgItems, v.data...)

	return pkgItems
}
