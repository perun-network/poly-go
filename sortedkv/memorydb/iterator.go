// SPDX-License-Identifier: Apache-2.0

package memorydb

// Iterator provides an iterator over a key range.
type Iterator struct {
	next   int
	keys   []string
	values []string
}

// Next returns true if the iterator has a next element.
func (i *Iterator) Next() bool {
	i.next++
	return i.next <= len(i.keys)
}

// Key returns the key of the current element.
func (i *Iterator) Key() string {
	if i.next == 0 {
		panic("Iterator.Key() accessed before Next() or after Close().")
	}

	if i.next > len(i.keys) {
		return ""
	}
	return i.keys[i.next-1]
}

// Value returns the value of the current element.
func (i *Iterator) Value() string {
	if i.next == 0 {
		panic("Iterator.Value() accessed before Next() or after Close().")
	}

	if i.next > len(i.keys) {
		return ""
	}
	return i.values[i.next-1]
}

// ValueBytes returns the value converted to bytes of the current element.
func (i *Iterator) ValueBytes() []byte {
	return []byte(i.Value())
}

// Close closes this iterator.
func (i *Iterator) Close() error {
	i.next = 0
	i.keys = nil
	i.values = nil
	return nil
}
