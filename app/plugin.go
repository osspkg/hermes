package app

import (
	"github.com/osspkg/goppy/plugins"
	"github.com/osspkg/hermes/app/addons"
)

var Plugins = plugins.Plugins{}.Inject(
	addons.Plugin,
)
