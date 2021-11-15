// SPDX-License-Identifier: Apache-2.0

package key_test

import (
	"testing"

	"polycry.pt/poly-go/sortedkv/key"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	assert.Equal(t, key.Next(""), "\x00")
	assert.Equal(t, key.Next("a"), "a\x00")
}

func TestIncPrefix(t *testing.T) {
	assert.Equal(t, key.IncPrefix(""), "")
	assert.Equal(t, key.IncPrefix("\x00"), "\x01")
	assert.Equal(t, key.IncPrefix("a"), "b")
	assert.Equal(t, key.IncPrefix("zoo"), "zop")
	assert.Equal(t, key.IncPrefix("\xff"), "")
	assert.Equal(t, key.IncPrefix("\xffa"), "\xffb")
	assert.Equal(t, key.IncPrefix("a\xff"), "b")
	assert.Equal(t, key.IncPrefix("\xff\xff\xff"), "")
}
