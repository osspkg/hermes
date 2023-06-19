package addons

import "github.com/osspkg/goppy/plugins"

var Plugin = plugins.Plugin{
	Config: &Config{},
	Inject: New,
}

type Config struct {
	Addons string `yaml:"addons"`
}

func (v *Config) Default() {
	v.Addons = "./addons"
}
