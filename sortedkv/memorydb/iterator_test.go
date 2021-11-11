// SPDX-License-Identifier: Apache-2.0

package memorydb

import (
	"testing"

	"polycry.pt/poly-go/sortedkv"
	"polycry.pt/poly-go/sortedkv/test"
)

func TestIterator(t *testing.T) {
	t.Run("Generic iterator test", func(t *testing.T) {
		test.GenericIteratorTest(t, NewDatabase())
	})

	t.Run("Table iterator test", func(t *testing.T) {
		test.GenericIteratorTest(t, sortedkv.NewTable(NewDatabase(), "table"))
	})
}
