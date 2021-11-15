// SPDX-License-Identifier: Apache-2.0

// Package atomic contains extensions of "sync/atomic"
package atomic // import "polycry.pt/poly-go/sync/atomic"

import "sync/atomic"

// Bool is an atomically accessible boolean. Its initial state is false.
type Bool int32

// IsSet returns whether the bool is set.
func (b *Bool) IsSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }

// Set atomically sets the bool to true.
func (b *Bool) Set() { atomic.StoreInt32((*int32)(b), 1) }

// TrySet atomically sets the bool to true and returns whether it was false
// before.
func (b *Bool) TrySet() bool { return atomic.SwapInt32((*int32)(b), 1) == 0 }

// Unset atomically sets the bool to false.
func (b *Bool) Unset() { atomic.StoreInt32((*int32)(b), 0) }

// TryUnset atomically sets the bool to false and returns whether it was true
// before.
func (b *Bool) TryUnset() bool { return atomic.SwapInt32((*int32)(b), 0) == 1 }
