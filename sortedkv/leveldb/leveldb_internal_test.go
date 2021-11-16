// SPDX-License-Identifier: Apache-2.0

package leveldb

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"polycry.pt/poly-go/sortedkv"
	"polycry.pt/poly-go/sortedkv/test"
)

func TestBatch(t *testing.T) {
	runTestOnTempDatabase(t, func(db *Database) {
		test.GenericBatchTest(t, db)
	})
	runTestOnTempDatabase(t, func(db *Database) {
		test.GenericBatchTest(t, sortedkv.NewTable(db, "table"))
	})
}

func TestDatabase(t *testing.T) {
	runTestOnTempDatabase(t, func(db *Database) {
		test.GenericDatabaseTest(t, db)
	})
}

func TestIterator(t *testing.T) {
	runTestOnTempDatabase(t, func(db *Database) {
		test.GenericIteratorTest(t, db)
	})
}

func runTestOnTempDatabase(t *testing.T, tester func(*Database)) {
	t.Helper()
	// Create a temporary directory and delete it when done
	path, err := ioutil.TempDir("", "poly_testdb_")
	require.Nil(t, err, "Could not create temporary directory for database")
	defer func() { require.Nil(t, os.RemoveAll(path)) }()

	// Create a database in the directory and close it when done
	db, err := LoadDatabase(path)
	require.Nil(t, err, "Could not load database")
	defer func() { assert.Nil(t, db.DB.Close(), "Could not close database") }()

	tester(db)
}
