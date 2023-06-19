package resolver

import "github.com/osspkg/goppy/plugins"

var Plugin = plugins.Plugin{
	Inject: New,
}
