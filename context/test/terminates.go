// Copyright 2020 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
