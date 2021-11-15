// SPDX-License-Identifier: Apache-2.0

package sync_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"polycry.pt/poly-go/context/test"
	"polycry.pt/poly-go/sync"
)

func TestWaitGoup(t *testing.T) {
	var wg sync.WaitGroup

	// Empty waitgroups have a counter of 0, and should immediately return.
	test.AssertTerminatesQuickly(t, wg.Wait)

	t.Run("negative counter", func(t *testing.T) {
		assert.Panics(t, func() { wg.Add(-5) })
		assert.Panics(t, wg.Done)
	})

	wg.Add(1)
	test.AssertNotTerminatesQuickly(t, wg.Wait)
	wg.Done()
	test.AssertTerminatesQuickly(t, wg.Wait)
	test.AssertTerminatesQuickly(t, wg.Wait)

	const N = 5
	wg.Add(N)
	test.AssertNotTerminatesQuickly(t, wg.Wait)

	for i := 0; i < N-1; i++ {
		wg.Done()
		test.AssertNotTerminatesQuickly(t, wg.Wait)
	}
	wg.Done()
	test.AssertTerminatesQuickly(t, wg.Wait)
}
