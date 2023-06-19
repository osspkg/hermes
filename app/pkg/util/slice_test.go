package util_test

import (
	"github.com/osspkg/hermes/app/pkg/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnit_FlipStringSlice(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e"}
	util.FlipStringSlice(data)
	require.Equal(t, []string{"e", "d", "c", "b", "a"}, data)

	data = []string{"a", "b", "c", "d"}
	util.FlipStringSlice(data)
	require.Equal(t, []string{"d", "c", "b", "a"}, data)
}
