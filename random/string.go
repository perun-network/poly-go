// SPDX-License-Identifier: Apache-2.0

package random

import (
	"math/rand"
)

var runes = []rune("abcdefghijklmnopABCDEFGHIJKLMNOP0123456789 ")

// String returns a random string of the given size with potentially lower and
// uppercase letters as well as intermingled numbers and spaces.
func String(rng *rand.Rand, size int) string {
	s := make([]rune, size)
	for i := range s {
		s[i] = runes[rng.Int()%len(runes)]
	}
	return string(s)
}
