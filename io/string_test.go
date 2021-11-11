// SPDX-License-Identifier: Apache-2.0

package io

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	polytest "polycry.pt/poly-go/test"
)

func TestEncodeDecodeString(t *testing.T) {
	assert := assert.New(t)
	rng := polytest.Prng(t)
	uint8buf, uint16buf := make([]byte, math.MaxUint8), make([]byte, math.MaxUint16)
	rng.Read(uint8buf)
	rng.Read(uint16buf)

	t.Run("valid strings", func(t *testing.T) {
		ss := []string{"", "a", "perun", string(uint8buf), string(uint16buf)}

		for _, s := range ss {
			r, w := io.Pipe()
			done := make(chan struct{})
			go func() {
				defer close(done)
				defer w.Close()
				assert.NoError(encodeString(w, s))
			}()

			var d string
			assert.NoError(decodeString(r, &d))
			r.Close()
			assert.Equal(s, d)
			<-done
		}
	})

	t.Run("too long string", func(t *testing.T) {
		tooLong := string(append(uint16buf, 42))
		var buf bytes.Buffer
		assert.Error(encodeString(&buf, tooLong))
		assert.Zero(buf.Len(), "nothing should have been written to the stream")
	})

	t.Run("short stream", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, byteOrder, uint16(16))
		buf.Write(make([]byte, 8)) // 8 bytes missing

		var d string
		assert.Error(decodeString(&buf, &d))
		assert.Zero(buf.Len(), "buffer should be exhausted")
	})
}
