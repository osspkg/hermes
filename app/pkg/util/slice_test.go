/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a LGPL-3.0 license that can be found in the LICENSE file.
 */

package util_test

import (
	"testing"

	"github.com/osspkg/hermes/app/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestUnit_FlipStringSlice(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e"}
	util.FlipStringSlice(data)
	require.Equal(t, []string{"e", "d", "c", "b", "a"}, data)

	data = []string{"a", "b", "c", "d"}
	util.FlipStringSlice(data)
	require.Equal(t, []string{"d", "c", "b", "a"}, data)
}
