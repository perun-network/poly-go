// SPDX-License-Identifier: Apache-2.0

// Package test tests the helper utilities regarding go contexts.
package test // import "polycry.pt/poly-go/context/test"

import (
	"context"
	"time"

	polyctx "polycry.pt/poly-go/context"
	"polycry.pt/poly-go/test"
)

// AssertTerminatesCtx asserts that a function terminates before a context is
// done.
func AssertTerminatesCtx(ctx context.Context, t test.T, fn func()) {
	t.Helper()

	if !polyctx.TerminatesCtx(ctx, fn) {
		t.Errorf("function should have terminated within timeout")
	}
}

// AssertNotTerminatesCtx asserts that a function does not terminate before a
// context is done.
func AssertNotTerminatesCtx(ctx context.Context, t test.T, fn func()) {
	t.Helper()

	if polyctx.TerminatesCtx(ctx, fn) {
		t.Errorf("Function should not have terminated within timeout")
	}
}

// AssertTerminates asserts that a function terminates within a certain
// timeout.
func AssertTerminates(t test.T, timeout time.Duration, fn func()) {
	t.Helper()

	if !polyctx.Terminates(timeout, fn) {
		t.Errorf("Function should have terminated within timeout")
	}
}

// AssertNotTerminates asserts that a function does not terminate within a
// certain timeout.
func AssertNotTerminates(t test.T, timeout time.Duration, fn func()) {
	t.Helper()

	if polyctx.Terminates(timeout, fn) {
		t.Errorf("Function should not have terminated within timeout")
	}
}

// AssertTerminatesQuickly asserts that a function terminates within 20 ms.
func AssertTerminatesQuickly(t test.T, fn func()) {
	t.Helper()

	if !polyctx.TerminatesQuickly(fn) {
		t.Errorf("Function should have terminated within timeout")
	}
}

// AssertNotTerminatesQuickly asserts that a function does not terminate within
// 20 ms.
func AssertNotTerminatesQuickly(t test.T, fn func()) {
	t.Helper()

	if polyctx.TerminatesQuickly(fn) {
		t.Errorf("Function should not have terminated within timeout")
	}
}
