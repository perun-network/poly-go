// SPDX-License-Identifier: Apache-2.0

package test

import (
	"runtime"
	"sync"
)

// Skipper is a subset of the testing.T functionality needed by OnlyOnce().
type Skipper interface {
	SkipNow()
}

var (
	executedTests      = make(map[string]struct{})
	executedTestsMutex sync.Mutex
)

// OnlyOnce records a test case and skips it if it already executed once.
// Test case identification is done by observing the stack. Calls SkipNow() on
// tests that have already been executed. OnlyOnce() has to be called directly
// from the test's function, as its first action.
func OnlyOnce(t Skipper) {
	// nolint:dogsled
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	executedTestsMutex.Lock()
	defer executedTestsMutex.Unlock()

	if _, executed := executedTests[name]; executed {
		t.SkipNow()
	}
	executedTests[name] = struct{}{}
}
