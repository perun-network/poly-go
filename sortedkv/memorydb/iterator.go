// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of go-perun. Use
// of this source code is governed by the Apache 2.0 license that can be found
// in the LICENSE file.

package memorydb

import "perun.network/go-perun/log"

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
		log.Panic("Iterator.Key() accessed before Next() or after Close().")
	}

	if i.next > len(i.keys) {
		return ""
	}
	return i.keys[i.next-1]
}

// Value returns the value of the current element.
func (i *Iterator) Value() string {
	if i.next == 0 {
		log.Panic("Iterator.Value() accessed before Next() or after Close().")
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
