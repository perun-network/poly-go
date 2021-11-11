// SPDX-License-Identifier: Apache-2.0

package context_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	polyctx "polycry.pt/poly-go/context"
)

func TestIsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.True(t, polyctx.IsContextError(errors.WithStack(ctx.Err())))

	// context that immediately times out
	ctx, cancel = context.WithTimeout(context.Background(), 0)
	defer cancel()
	assert.True(t, polyctx.IsContextError(errors.WithStack(ctx.Err())))

	assert.False(t, polyctx.IsContextError(errors.New("no context error")))
}

func TestIsDone(t *testing.T) {
	assert.False(t, polyctx.IsDone(context.Background()))

	ctx, cancel := context.WithCancel(context.Background())
	assert.False(t, polyctx.IsDone(ctx))
	cancel()
	assert.True(t, polyctx.IsDone(ctx))

	// context that immediately times out
	ctx, cancel = context.WithTimeout(context.Background(), 0)
	defer cancel()
	assert.True(t, polyctx.IsDone(ctx))
}
