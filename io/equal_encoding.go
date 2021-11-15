// SPDX-License-Identifier: Apache-2.0

package io

import (
	"bytes"

	"github.com/pkg/errors"
)

// EqualEncoding returns whether the two Encoders `a` and `b` encode to the same byteslice
// or an error when the encoding failed.
func EqualEncoding(a, b Encoder) (bool, error) {
	buffA := new(bytes.Buffer)
	buffB := new(bytes.Buffer)

	// golang does not have a XOR
	if (a == nil) != (b == nil) {
		return false, errors.New("only one argument was nil")
	}
	// just using a == b would be too easy here since go panics
	if (a == nil) && (b == nil) {
		return true, nil
	}

	if err := a.Encode(buffA); err != nil {
		return false, errors.Wrap(err, "EqualEncoding encode error")
	}
	if err := b.Encode(buffB); err != nil {
		return false, errors.Wrap(err, "EqualEncoding encode error")
	}

	return bytes.Equal(buffA.Bytes(), buffB.Bytes()), nil
}
