// SPDX-License-Identifier: Apache-2.0

// Package sync contains a mutex that can be used in a select statement.
package sync // import "polycry.pt/poly-go/sync"

import (
	"context"
	"sync"
)

// Mutex is a replacement of the standard mutex type.
// It supports the additional TryLock() function, as well as a variant that can
// be used in a select statement.
type Mutex struct {
	locked chan struct{} // The internal mutex is modelled by a channel.
	once   sync.Once     // Needed to initialize the mutex on its first use.
}

// initOnce initialises the mutex if it has not been initialised yet.
func (m *Mutex) initOnce() {
	m.once.Do(func() { m.locked = make(chan struct{}, 1) })
}

// Lock blockingly locks the mutex.
func (m *Mutex) Lock() {
	m.initOnce()
	m.locked <- struct{}{}
}

// TryLock tries to lock the mutex without blocking.
// Returns whether the mutex was acquired.
func (m *Mutex) TryLock() bool {
	m.initOnce()
	select {
	case m.locked <- struct{}{}:
		return true
	default:
		return false
	}
}

// TryLockCtx tries to lock the mutex within a timeout provided by a context.
// For an instant timeout, a nil context has to be passed. Returns whether the
// mutex was acquired.
func (m *Mutex) TryLockCtx(ctx context.Context) bool {
	m.initOnce()

	if ctx == nil {
		return m.TryLock()
	}

	// Check for the deadline first because 'select' choses a random
	// available case.
	select {
	case <-ctx.Done():
		return false
	default:
	}

	select {
	case m.locked <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

// Unlock unlocks the mutex.
// If the mutex was not locked, panics.
func (m *Mutex) Unlock() {
	select {
	case <-m.locked:
	default:
		panic("tried to unlock unlocked mutex")
	}
}
