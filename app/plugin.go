package app

import (
	"github.com/osspkg/goppy/plugins"
	"github.com/osspkg/hermes/app/addons"
	"github.com/osspkg/hermes/app/resolver"
)

var Plugins = plugins.Plugins{}.Inject(
	addons.Plugin,
	resolver.Plugin,
)
