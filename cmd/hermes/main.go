/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a LGPL-3.0 license that can be found in the LICENSE file.
 */

package main

import (
	"github.com/osspkg/goppy"
	"github.com/osspkg/goppy/plugins/database"
	"github.com/osspkg/goppy/plugins/web"
	hermes "github.com/osspkg/hermes/app"
)

func main() {
	app := goppy.New()
	app.WithConfig("./config.yaml") // Reassigned via the `--config` argument when run via the console.
	app.Plugins(
		web.WithHTTP(),
		database.WithMySQL(),
	)
	app.Plugins(hermes.Plugins...)
	app.Run()
}
