// SPDX-License-Identifier: Apache-2.0

package test

import (
	"time"
)

// eventuallyT is an implementation of T that records whether any Error or Fatal
// call has been made on it.
type eventuallyT struct {
	called bool
}

func (et *eventuallyT) Errorf(string, ...interface{}) {
	et.fail()
}

func (et *eventuallyT) FailNow() {
	et.fail()
}

func (et *eventuallyT) fail() {
	et.called = true
	panic("eventually")
}

func (et *eventuallyT) Helper() {}

// An EventuallyTest has fixed `within` and `pause` duration parameters so that
// a test for eventual success can be run by just passing the testing object and
// the test function to method Eventually.
type EventuallyTest struct {
	within, pause time.Duration
}

// NewEventually creates a new EventuallyTest object which fixes the `within`
// and `pause` duration parameters.
func NewEventually(within, pause time.Duration) *EventuallyTest {
	return &EventuallyTest{
		within: within,
		pause:  pause,
	}
}

// Eventually does the same as the free function of the same name but with the
// duration parameters `within` and `pause` taken from the EventuallyTest `et`.
func (et *EventuallyTest) Eventually(t T, assertion func(T)) {
	until := time.Now().Add(et.within)
	for time.Now().Before(until) {
		if et.test(assertion) {
			return
		}

		if remaining := time.Until(until); remaining < et.pause {
			time.Sleep(remaining)
		} else {
			time.Sleep(et.pause)
		}
	}

	assertion(t) // final call with actual testing object.
}

func (et *EventuallyTest) test(assertion func(T)) (success bool) {
	e := new(eventuallyT)

	success = false
	defer func() {
		if !success && e.called {
			_ = recover()
		}
	}()

	assertion(e)
	return true
}

// Eventually runs the test `test` until it stops failing for the duration
// `within` while sleeping for `pause` in between test executions. The final
// call to `test` with the actual test object `t` is guaranteed to be run
// exactly at time time.Now().Add(within).
//
// The test should be a read-only test that can safely be run several times
// without changing the tested objects. Until the final test call, any call to
// Error or Fail aborts execution of the test function by panicking, to avoid
// running unnecessary checks.
//
// Eventually does not start any go routines.
func Eventually(t T, test func(T), within, pause time.Duration) {
	NewEventually(within, pause).Eventually(t, test)
}

const (
	within100msTimeout = 100 * time.Millisecond
	within100msPause   = 5 * time.Millisecond
	within1sTimeout    = 1 * time.Second
	within1sPause      = 20 * time.Millisecond
	within10sTimeout   = 10 * time.Second
	within10sPause     = 200 * time.Millisecond
)

// Within100ms is an EventuallyTest that runs the test every 5ms up to 100ms.
var Within100ms = NewEventually(within100msTimeout, within100msPause)

// Within1s is an EventuallyTest that runs the test every 20ms up to 1s.
var Within1s = NewEventually(within1sTimeout, within1sPause)

// Within10s is an EventuallyTest that runs the test every 200ms up to 10s.
var Within10s = NewEventually(within10sTimeout, within10sPause)
