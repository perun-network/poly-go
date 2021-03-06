// SPDX-License-Identifier: Apache-2.0

package sync_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"polycry.pt/poly-go/context/test"
	"polycry.pt/poly-go/sync"
)

const timeout = 100 * time.Millisecond

func TestCloser_Closed(t *testing.T) {
	t.Parallel()
	var c sync.Closer

	assert.NotNil(t, c.Closed())
	select {
	case _, ok := <-c.Closed():
		t.Fatalf("Closed() should not yield a value, ok = %t", ok)
	default:
	}

	require.NoError(t, c.Close())

	test.AssertTerminates(t, timeout, func() {
		_, ok := <-c.Closed()
		assert.False(t, ok)
	})
}

func TestCloser_Ctx(t *testing.T) {
	t.Parallel()
	var c sync.Closer
	ctx := c.Ctx()
	assert.NoError(t, ctx.Err())
	assert.Nil(t, ctx.Value(nil))
	_, ok := ctx.Deadline()
	assert.False(t, ok)

	select {
	case <-ctx.Done():
		t.Error("context should not be closed")
	default: // expected
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		<-ctx.Done()
		assert.Same(t, ctx.Err(), context.Canceled)
	}()
	assert.NoError(t, c.Close())
	<-done
}
