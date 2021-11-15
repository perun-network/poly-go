// SPDX-License-Identifier: Apache-2.0

package sync_test

import (
	"context"
	stdsync "sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"polycry.pt/poly-go/context/test"
	"polycry.pt/poly-go/sync"
)

func TestSignal_Signal(t *testing.T) {
	const N = 64

	s := sync.NewSignal()

	var done stdsync.WaitGroup
	waitNTimes(s, &done, N)

	for i := 0; i < N; i++ {
		go s.Signal()
	}

	test.AssertTerminatesQuickly(t, done.Wait)
}

func TestSignal_Broadcast(t *testing.T) {
	const N = 20

	s := sync.NewSignal()

	for x := 0; x < 4; x++ {
		var done stdsync.WaitGroup
		waitNTimes(s, &done, N)

		// ensure that all goroutines are at s.Wait().
		time.Sleep(100 * time.Millisecond)

		s.Broadcast()
		test.AssertTerminatesQuickly(t, done.Wait)
	}
}

func waitNTimes(s *sync.Signal, wg *stdsync.WaitGroup, n int) {
	started := make(chan struct{})
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			started <- struct{}{}
			defer wg.Done()
			s.Wait()
		}()
		// Ensure the coroutine is already started to prevent scheduling delays
		// on slow pipelines.
		<-started
	}
}

func TestSignal_Wait(t *testing.T) {
	s := sync.NewSignal()
	test.AssertNotTerminatesQuickly(t, s.Wait)

	go func() {
		time.Sleep(200 * time.Millisecond)
		s.Broadcast()
	}()
	test.AssertNotTerminates(t, 100*time.Millisecond, s.Wait)
	test.AssertTerminates(t, 200*time.Millisecond, s.Wait)
}

func TestSignal_WaitCtx(t *testing.T) {
	s := sync.NewSignal()

	timeout := 100 * time.Millisecond
	go func() {
		time.Sleep(3 * timeout)
		s.Signal()
	}()
	test.AssertNotTerminates(t, timeout, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*timeout)
		defer cancel()
		assert.False(t, s.WaitCtx(ctx))
	})
	test.AssertTerminates(t, 3*timeout, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*timeout)
		defer cancel()
		assert.True(t, s.WaitCtx(ctx))
	})
}
