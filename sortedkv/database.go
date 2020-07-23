// Copyright 2019 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package sortedkv defines a sorted key-value store abstraction. It is used by
// other persistence packages (like channel/persistence) to continuously save
// the state of the running program or load it upon startup. It is based on
// go-ethereum's https://github.com/ethereum/go-ethereum/ethdb and
// https://github.com/ethereum/go-ethereum/core/rawdb and is also inspired by
// perkeep.org/pkg/sorted
package sortedkv // import "perun.network/go-perun/pkg/sortedkv"

import "io"

// ErrNotFound is returned whenever a key is not in the db.
type ErrNotFound struct {
	Key string
}

// Error returns the error string.
func (e *ErrNotFound) Error() string {
	return "sortedkv: key not found: " + e.Key
}

// Reader wraps the Has, Get and GetBytes methods of a key-value store.
type Reader interface {
	// Has checks if a key is present in the store.
	Has(key string) (bool, error)

	// Get returns the value as string for given key if it is present in the store.
	Get(key string) (string, error)

	// GetBytes returns the value as []byte for given key if it is present in the store.
	GetBytes(key string) ([]byte, error)
}

// Writer wraps the Put and Delete methods of a key-value store.
type Writer interface {
	// Put inserts the given value into the key-value store.
	// If the key is already present, it is overwritten and no error is returned.
	Put(key string, value string) error

	// PutBytes inserts the given value into the key-value store.
	// If the key is already present, it is overwritten and no error is returned.
	PutBytes(key string, value []byte) error

	// Delete removes the key from the key-value store.
	// If the key is not present, an error is returned
	Delete(key string) error
}

// Database is a key-value store (not to be confused with SQL-like databases).
type Database interface {
	Reader
	Writer
	Batcher
	Iterable
	io.Closer
}
