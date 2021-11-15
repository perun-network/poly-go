// SPDX-License-Identifier: Apache-2.0

package sortedkv

// tableBatch is a wrapper around a Database Batch with a key prefix. All
// Writer operations are automatically prefixed.
type tableBatch struct {
	Batch
	prefix string
}

func (b *tableBatch) pkey(key string) string {
	return b.prefix + key
}

// Put puts a value into a table batch.
func (b *tableBatch) Put(key, value string) error {
	return b.Batch.Put(b.pkey(key), value)
}

// Put puts a value into a table batch.
func (b *tableBatch) PutBytes(key string, value []byte) error {
	return b.Batch.PutBytes(b.pkey(key), value)
}

// Delete deletes a value from a table batch.
func (b *tableBatch) Delete(key string) error {
	return b.Batch.Delete(b.pkey(key))
}
