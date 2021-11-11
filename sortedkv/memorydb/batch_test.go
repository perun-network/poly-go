// SPDX-License-Identifier: Apache-2.0

package memorydb

import (
	"testing"

	"polycry.pt/poly-go/sortedkv"
	"polycry.pt/poly-go/sortedkv/test"
)

func TestBatch(t *testing.T) {
	t.Run("Generic Batch test", func(t *testing.T) {
		test.GenericBatchTest(t, NewDatabase())
	})

	t.Run("Generic table batch test", func(t *testing.T) {
		test.GenericBatchTest(t, sortedkv.NewTable(NewDatabase(), "table"))
	})
}
