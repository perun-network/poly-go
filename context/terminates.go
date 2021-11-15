// SPDX-License-Identifier: Apache-2.0

package context

import (
	"context"
	"time"
)

// TerminatesCtx checks whether a function terminates before a context is done.
func TerminatesCtx(ctx context.Context, fn func()) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	done := make(chan struct{}, 1)
	go func() {
		fn()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return true
	case <-ctx.Done():
		return false
	}
}

// Terminates checks whether a function terminates within a certain timeout.
func Terminates(timeout time.Duration, fn func()) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return TerminatesCtx(ctx, fn)
}

// TerminatesQuickly checks whether a function terminates within 20 ms.
func TerminatesQuickly(fn func()) bool {
	return Terminates(time.Millisecond*20, fn)
}
