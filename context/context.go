// SPDX-License-Identifier: Apache-2.0

// Package context contains helper utilities regarding go contexts.
package context // import "polycry.pt/poly-go/context"

import (
	"context"

	"github.com/pkg/errors"
)

// IsContextError returns whether the given error originates from a context that
// was cancelled or whose deadline exceeded. Prior to checking, the error is
// unwrapped by calling errors.Cause.
func IsContextError(err error) bool {
	err = errors.Cause(err)
	return err == context.Canceled || err == context.DeadlineExceeded
}

// IsDone returns whether ctx is done.
func IsDone(ctx interface{ Done() <-chan struct{} }) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
