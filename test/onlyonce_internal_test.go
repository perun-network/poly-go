// SPDX-License-Identifier: Apache-2.0

package test

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

var onlyOnceTestCalls int32

func TestOnlyOnce1(t *testing.T) {
	testOnlyOnce(t)
}

func TestOnlyOnce2(t *testing.T) {
	testOnlyOnce(t)
}

func testOnlyOnce(t *testing.T) {
	OnlyOnce(t)
	assert.Equal(t, int32(1), atomic.AddInt32(&onlyOnceTestCalls, 1))
}
