// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"testing"
	"time"

	"polycry.pt/poly-go/test"
)

const timeout = 200 * time.Millisecond

func TestAssertTerminatesCtx(t *testing.T) {
	t.Run("error case", func(t *testing.T) {
		test.AssertError(t, func(t test.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			AssertTerminatesCtx(ctx, t, func() {})
		})
	})

	t.Run("success case", func(t *testing.T) {
		AssertTerminatesCtx(context.Background(), t, func() {})
	})
}

func TestAssertNotTerminatesCtx(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		AssertNotTerminatesCtx(ctx, t, func() {})
	})

	t.Run("error case", func(t *testing.T) {
		test.AssertError(t, func(t test.T) {
			AssertNotTerminatesCtx(context.Background(), t, func() {})
		})
	})
}

func TestAssertTerminates(t *testing.T) {
	t.Run("error case", func(t *testing.T) {
		test.AssertError(t, func(t test.T) {
			AssertTerminates(t, -1, func() {})
		})
	})

	t.Run("success case", func(t *testing.T) {
		AssertTerminates(t, timeout, func() {})
	})
}

func TestAssertNotTerminates(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		AssertNotTerminates(t, -1, func() {})
	})

	t.Run("error case", func(t *testing.T) {
		test.AssertError(t, func(t test.T) {
			AssertNotTerminates(t, timeout, func() {})
		})
	})
}

func TestAssertTerminatesQuickly(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		AssertTerminatesQuickly(t, func() {})
	})

	t.Run("error case", func(t *testing.T) {
		test.AssertError(t, func(t test.T) {
			AssertTerminatesQuickly(t, func() { time.Sleep(time.Hour) })
		})
	})
}

func TestAssertNotTerminatesQuickly(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		AssertNotTerminatesQuickly(t, func() { time.Sleep(time.Hour) })
	})

	t.Run("error case", func(t *testing.T) {
		test.AssertError(t, func(t test.T) {
			AssertNotTerminatesQuickly(t, func() {})
		})
	})
}
