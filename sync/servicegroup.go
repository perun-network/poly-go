// SPDX-License-Identifier: Apache-2.0

package sync

import (
	"sync"
	"sync/atomic"
)

// A ServiceGroup is a collection of subservices that comprise an overall
// service. It is intended to be used in cases where the overall service is
// considered failed from the moment a single subservice fails.
type ServiceGroup struct {
	n       int32
	finOnce sync.Once

	done chan struct{}
	err  error
}

// NewServiceGroup returns a new ServiceGroup.
func NewServiceGroup() *ServiceGroup {
	return &ServiceGroup{done: make(chan struct{})}
}

// Go starts f in a new goroutine.
//
// The channel returned by Done is closed once the first goroutine returns a
// non-nil error or all goroutines return with nil.
//
// Note that if the first goroutine started with Go returns a nil error before
// any second goroutine is started, the ServiceGroup will already be done with a
// nil error that will stay nil even if the second goroutine returns a non-nil
// error. This is inteded, as ServiceGroup is an eagerly-finishing version of
// golang.org/x/sync/errgroup.Group.
func (g *ServiceGroup) Go(f func() error) {
	atomic.AddInt32(&g.n, 1)

	go func() {
		var err error
		// Make sure that the ServiceGroup finishes even when f panicks.
		defer func() {
			if rem := atomic.AddInt32(&g.n, -1); err == nil && rem > 0 {
				return
			}
			// Now this is either a failed routine or the last.
			g.finish(err)
		}()
		err = f()
	}()
}

func (g *ServiceGroup) finish(err error) {
	g.finOnce.Do(func() {
		g.err = err
		close(g.done)
	})
}

// Done returns a signalling channel that is closed once either the first
// routine started with Go returns a non-nil error or all goroutines have
// returned a nil error.
func (g *ServiceGroup) Done() <-chan struct{} { return g.done }

// Err is set to the error of the first failing routine that was started with
// Go, if any. It returns nil before the channel returned by Done is closed.
func (g *ServiceGroup) Err() error { return g.err }

// Wait is a shortcut for waiting on Done() and then returning Err().
func (g *ServiceGroup) Wait() error {
	<-g.done
	return g.err
}
