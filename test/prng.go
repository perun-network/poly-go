// SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	envTestSeed          = "GOTESTSEED"
	genRootSeedPrintBase = 10
	genRootSeedPrintLen  = 64
)

var rootSeed int64

func init() {
	rootSeed = genRootSeed()
	fmt.Printf("pkg/test: using rootSeed %d\n", rootSeed) // nolint: forbidigo
}

func genRootSeed() (rootSeed int64) {
	if val, ok := os.LookupEnv(envTestSeed); ok {
		rootSeed, err := strconv.ParseInt(val, genRootSeedPrintBase, genRootSeedPrintLen)
		if err != nil {
			panic(fmt.Sprintf("Could not parse %s = '%s' as int64", envTestSeed, val))
		}
		return rootSeed
	}
	return time.Now().UnixNano()
}

// Prng returns a pseudo-RNG that is seeded with the output of the `Seed`
// function by passing it `t.Name()`.
// Use it in tests with: rng := pkgtest.Prng(t).
func Prng(t interface{ Name() string }, args ...interface{}) *rand.Rand {
	return rand.New(rand.NewSource(Seed(t.Name(), args...))) // nolint: gosec
}

// Seed generates a seed that is dependent on the rootSeed and the passed
// seed argument.
// To fix this seed, set the GOTESTSEED environment variable.
// Example: GOTESTSEED=123 go test ./...
// Does not work with function pointers or structs without public fields.
func Seed(seed string, args ...interface{}) int64 {
	hasher := fnv.New64a()
	enc := gob.NewEncoder(hasher)
	if err := enc.Encode(seed); err != nil {
		panic("Could not gob-encode seed")
	}
	for _, arg := range args {
		if err := enc.Encode(arg); err != nil {
			panic(fmt.Sprintf("Could not gob-encode value: %v", err))
		}
	}
	if err := binary.Write(hasher, binary.LittleEndian, rootSeed); err != nil {
		panic("Could not hash the root seed")
	}
	return int64(hasher.Sum64())
}

// NameStr endows a string with a function `Name() string` that returns itself.
// This is helpful for using strings as a first argument to `Prng`.
type NameStr string

// Name returns the string underlying NameStr.
func (s NameStr) Name() string { return string(s) }
