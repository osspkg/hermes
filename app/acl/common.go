package acl

import "sync"

type ACL struct {
	data map[string]map[uint]uint
	mux  sync.RWMutex
}

func New() *ACL {
	return &ACL{
		data: make(map[string]map[uint]uint, 100),
	}
}

func (v *ACL) Setup(pkgName string, id uint, forms []uint) {
	v.mux.Lock()
	defer v.mux.Unlock()

	list, ok := v.data[pkgName]
	if !ok {
		list = make(map[uint]uint, 20)
	}
	for _, form := range forms {
		list[form] = id
	}
	v.data[pkgName] = list
}

func (v *ACL) Clean(pkgName string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	delete(v.data, pkgName)
}
