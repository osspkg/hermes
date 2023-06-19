package main

import (
	"github.com/osspkg/goppy"
	"github.com/osspkg/goppy/plugins/web"
	hermes "github.com/osspkg/hermes/app"
)

func main() {
	app := goppy.New()
	app.WithConfig("./config.yaml") // Reassigned via the `--config` argument when run via the console.
	app.Plugins(
		web.WithHTTP(),
	)
	app.Plugins(hermes.Plugins...)
	app.Run()
}
