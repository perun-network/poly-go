// SPDX-License-Identifier: Apache-2.0

// Package test contains helper functions for testing.
package test

// CheckPanic tests whether a supplied function panics during its execution.
// Returns whether the supplied function did panic, and if so, also returns the
// value passed to panic().
func CheckPanic(function func()) (didPanic bool, value interface{}) {
	// Catch the panic, if it happens and store the passed value.
	defer func() {
		value = recover()
		if p, ok := value.(*Panic); ok {
			value = p.Value()
		}
	}()
	// Set up for the panic case.
	didPanic = true
	// Call the function to be checked.
	function()
	// This is only executed if panic() was not called.
	didPanic = false
	return
}
