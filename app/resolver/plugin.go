/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a LGPL-3.0 license that can be found in the LICENSE file.
 */

package resolver

import "github.com/osspkg/goppy/plugins"

var Plugin = plugins.Plugin{
	Inject: New,
}
