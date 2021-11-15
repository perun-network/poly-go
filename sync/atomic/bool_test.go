// SPDX-License-Identifier: Apache-2.0

package atomic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"polycry.pt/poly-go/sync/atomic"
)

func TestBool(t *testing.T) {
	assert := assert.New(t)

	var b atomic.Bool
	assert.False(b.IsSet())
	b.Set()
	assert.True(b.IsSet())
	assert.False(b.TrySet())
	assert.True(b.IsSet())

	b.Unset()
	assert.False(b.IsSet())
	assert.True(b.TrySet())
	assert.True(b.IsSet())
	assert.True(b.TryUnset())
	assert.False(b.TryUnset())
	assert.False(b.IsSet())
}
