// SPDX-License-Identifier: Apache-2.0

package test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAbort(t *testing.T) {
	abort := CheckAbort(func() {})
	assert.Nil(t, abort)

	abort = CheckAbort(func() { panic(1) })
	require.IsType(t, abort, (*Panic)(nil))
	assert.Equal(t, abort.(*Panic).Value(), 1)

	abort = CheckAbort(func() { panic(nil) })
	require.IsType(t, abort, (*Panic)(nil))
	assert.Nil(t, abort.(*Panic).Value())

	abort = CheckAbort(runtime.Goexit)
	assert.IsType(t, abort, (*Goexit)(nil))
}

func TestCheckGoexit(t *testing.T) {
	assert.True(t, CheckGoexit(runtime.Goexit))
	assert.Panics(t, func() { CheckGoexit(func() { panic("") }) })
	didPanic, pval := CheckPanic(func() { CheckGoexit(func() { panic(nil) }) })
	assert.True(t, didPanic)
	assert.Nil(t, pval)
	assert.False(t, CheckGoexit(func() {}))
}
