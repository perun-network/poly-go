// SPDX-License-Identifier: Apache-2.0

package io_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	ctxtest "polycry.pt/poly-go/context/test"
	polyio "polycry.pt/poly-go/io"
	iotest "polycry.pt/poly-go/io/test"
)

func TestByteSlice(t *testing.T) {
	var v1, v2, v3, v4 polyio.ByteSlice = []byte{}, []byte{255}, []byte{1, 2, 3, 4, 5, 6}, make([]byte, 20000)
	testByteSlices(t, v1, v2, v3, v4)
	iotest.GenericBrokenPipeTest(t, &v1, &v2, &v3, &v4)
}

// TestStutter tests what happens if the network stutters (split one message into several network packages).
func TestStutter(t *testing.T) {
	values := []byte{0, 1, 2, 3, 4, 5, 6, 255}
	r, w := io.Pipe()

	go func() {
		for _, v := range values {
			_, err := w.Write([]byte{v})
			assert.NoError(t, err)
		}
	}()

	var decodedValue polyio.ByteSlice = make([]byte, len(values))
	ctxtest.AssertTerminatesQuickly(t, func() {
		assert.NoError(t, decodedValue.Decode(r))
	})
	for i, v := range values {
		assert.Equal(t, decodedValue[i], v)
	}
}

func testByteSlices(t *testing.T, serial ...polyio.ByteSlice) {
	t.Helper()
	a := assert.New(t)
	r, w := io.Pipe()
	done := make(chan struct{})

	go func() {
		for _, v := range serial {
			a.NoError(v.Encode(w))
		}
		close(done)
	}()

	for i, v := range serial {
		d := make([]byte, len(v))
		dest := polyio.ByteSlice(d)

		a.NoError(dest.Decode(r), "failed to decode element")

		if !reflect.DeepEqual(v, dest) {
			t.Errorf("encoding and decoding the %dth element (%T) resulted in different value: %v, %v", i, v, reflect.ValueOf(v).Elem(), dest)
		}
	}
	<-done
}
