// SPDX-License-Identifier: Apache-2.0

package sortedkv

// Batch is a write-only database that buffers changes to the underlying
// database until Apply() is called.
type Batch interface {
	Writer // Put and Delete

	// Apply performs all batched actions on the database.
	Apply() error

	// Reset resets the batch so that it doesn't contain any items and can be reused.
	Reset()
}

// Batcher wraps the NewBatch method of a backing data store.
type Batcher interface {
	// NewBatch creates a Batch that will write to the Batcher.
	NewBatch() Batch
}
