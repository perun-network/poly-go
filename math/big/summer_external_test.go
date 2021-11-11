// SPDX-License-Identifier: Apache-2.0

package big_test

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	polybig "polycry.pt/poly-go/math/big"
	"polycry.pt/poly-go/test"
)

func TestSummer(t *testing.T) {
	rng := test.Prng(t)

	// Same summers do not error and are equal.
	for i := 0; i < 10; i++ {
		a := randomSummer(i, rng)
		equal, err := polybig.EqualSum(a, a)
		require.NoError(t, err)
		assert.True(t, equal)
	}

	// Same size summers do not error and are different.
	for i := 1; i < 10; i++ {
		a, b := randomSummer(i, rng), randomSummer(i, rng)
		equal, err := polybig.EqualSum(a, b)
		require.NoError(t, err)
		assert.False(t, equal)
	}

	// Different size summers do error.
	for i := 0; i < 10; i++ {
		a, b := randomSummer(i, rng), randomSummer(i+1, rng)
		equal, err := polybig.EqualSum(a, b)
		require.EqualError(t, err, "dimension mismatch")
		assert.False(t, equal)
	}
}

func TestAddSums(t *testing.T) {
	rng := test.Prng(t)

	// Calculates the correct sum.
	a, b := polybig.Sum{big.NewInt(4)}, polybig.Sum{big.NewInt(5)}
	sum, err := polybig.AddSums(a, b)
	require.NoError(t, err)
	require.Len(t, sum, 1)
	assert.Equal(t, sum[0], big.NewInt(9)) // 5 + 4 = 9

	// Different size summers do error.
	for i := 0; i < 10; i++ {
		a, b := randomSummer(i, rng), randomSummer(i+1, rng)
		_, err := polybig.AddSums(a, b)
		require.EqualError(t, err, "dimension mismatch")
	}

	// No input must return nil, nil.
	sum, err = polybig.AddSums()
	assert.Nil(t, err)
	assert.Nil(t, sum)
}

func randomSummer(len int, rng *rand.Rand) *polybig.Sum {
	data := make([]*big.Int, len)
	for i := range data {
		d, err := crand.Int(rng, new(big.Int).Lsh(big.NewInt(1), 255))
		if err != nil {
			panic(fmt.Sprintf("Creating random big.Int: %v", err))
		}
		data[i] = d
	}
	ret := polybig.Sum(data)
	return &ret
}

// TestAddSums_Const checks that `AddSums` does not modify the input.
func TestAddSums_Const(t *testing.T) {
	a, b := polybig.Sum{big.NewInt(4)}, polybig.Sum{big.NewInt(5)}

	_, err := polybig.AddSums(a, b)
	require.NoError(t, err)

	assert.Equal(t, a[0], big.NewInt(4))
	assert.Equal(t, b[0], big.NewInt(5))
}
