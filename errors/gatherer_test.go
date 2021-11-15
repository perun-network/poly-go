// SPDX-License-Identifier: Apache-2.0

package errors_test

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"polycry.pt/poly-go/context/test"
	"polycry.pt/poly-go/errors"
)

func TestGatherer_Failed(t *testing.T) {
	g := errors.NewGatherer()

	select {
	case <-g.Failed():
		t.Fatal("Failed must not be closed")
	default:
	}

	g.Add(stderrors.New(""))

	select {
	case <-g.Failed():
	default:
		t.Fatal("Failed must be closed")
	}
}

func TestGatherer_Go_and_Wait(t *testing.T) {
	g := errors.NewGatherer()

	const timeout = 100 * time.Millisecond

	g.Go(func() error {
		time.Sleep(timeout)
		return stderrors.New("")
	})

	test.AssertNotTerminates(t, timeout/2, func() { g.Wait() })
	var err error
	test.AssertTerminates(t, timeout, func() { err = g.Wait() })
	require.Error(t, err)
}

func TestGatherer_Go_and_DoneOrFailed(t *testing.T) {
	const timeout = 100 * time.Millisecond

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		g := errors.NewGatherer()

		// Two successful routines
		g.Go(func() error {
			time.Sleep(timeout / 2)
			return nil
		})
		g.Go(func() error {
			time.Sleep(timeout * 2)
			return nil
		})

		test.AssertNotTerminates(t, timeout, g.WaitDoneOrFailed)
		test.AssertTerminates(t, timeout*2, g.WaitDoneOrFailed)
	})
	t.Run("error", func(t *testing.T) {
		t.Parallel()
		g := errors.NewGatherer()

		// One slow successful and one error routine
		g.Go(func() error {
			time.Sleep(timeout * 10)
			return nil
		})
		g.Go(func() error {
			time.Sleep(timeout)
			return stderrors.New("")
		})

		test.AssertNotTerminates(t, timeout/2, g.WaitDoneOrFailed)
		test.AssertTerminates(t, timeout, g.WaitDoneOrFailed)
	})
	t.Run("ctx", func(t *testing.T) {
		t.Parallel()
		g := errors.NewGatherer()

		// One slow routine
		g.Go(func() error {
			time.Sleep(timeout * 10)
			return nil
		})

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		test.AssertNotTerminates(t, timeout/2, func() { g.WaitDoneOrFailedCtx(ctx) })
		test.AssertTerminates(t, timeout, func() { g.WaitDoneOrFailedCtx(ctx) })
	})
}

func TestGatherer_Add_and_Err(t *testing.T) {
	g := errors.NewGatherer()

	require.NoError(t, g.Err())

	g.Add(stderrors.New("1"))
	g.Add(stderrors.New("2"))
	require.Error(t, g.Err())
	require.Len(t, errors.Causes(g.Err()), 2)
}

func TestCauses(t *testing.T) {
	g := errors.NewGatherer()
	require.Len(t, errors.Causes(g.Err()), 0)

	g.Add(stderrors.New("1"))
	require.Len(t, errors.Causes(g.Err()), 1)

	g.Add(stderrors.New("2"))
	require.Len(t, errors.Causes(g.Err()), 2)

	g.Add(stderrors.New("3"))
	require.Len(t, errors.Causes(g.Err()), 3)

	err := stderrors.New("normal")
	causes := errors.Causes(err)
	require.Len(t, causes, 1)
	assert.Same(t, causes[0], err)
}

func TestAccumulatedError_Error(t *testing.T) {
	g := errors.NewGatherer()
	g.Add(stderrors.New("1"))
	require.Equal(t, g.Err().Error(), "(1 error)\n1): 1")

	g.Add(stderrors.New("2"))
	require.Equal(t, g.Err().Error(), "(2 errors)\n1): 1\n2): 2")
}

type stackTracer interface {
	StackTrace() pkgerrors.StackTrace
}

func TestAccumulatedError_StackTrace(t *testing.T) {
	g := errors.NewGatherer()

	g.Add(stderrors.New("1"))
	assert.Nil(t, g.Err().(stackTracer).StackTrace())
	g.Add(pkgerrors.New("2"))
	assert.NotNil(t, g.Err().(stackTracer).StackTrace())
}

func TestGatherer_OnFail(t *testing.T) {
	var (
		assert  = assert.New(t)
		g       = errors.NewGatherer()
		called  bool
		called2 bool
	)

	g.OnFail(func() {
		select {
		case <-g.Failed():
		case <-time.After(time.Second):
			assert.Fail("Failed not closed before OnFail hooks are executed")
		}

		called = true
		assert.False(called2)
	})

	g.OnFail(func() {
		assert.True(called)
		called2 = true
	})

	g.Add(nil)
	assert.False(called)
	assert.False(called2)

	g.Add(stderrors.New("error"))
	assert.True(called)
	assert.True(called2)
}
