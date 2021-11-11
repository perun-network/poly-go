// SPDX-License-Identifier: Apache-2.0

// Package test contains the generic serializer tests.
package test // import "polycry.pt/poly-go/io/test"

import (
	"io"
	"reflect"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	polyio "polycry.pt/poly-go/io"
)

// GenericSerializerTest runs multiple tests to check whether encoding
// and decoding of serializer values works.
func GenericSerializerTest(t *testing.T, serializers ...polyio.Serializer) {
	genericDecodeEncodeTest(t, serializers...)
	GenericBrokenPipeTest(t, serializers...)
}

// genericDecodeEncodeTest tests whether encoding and then decoding
// serializer values results in the original values.
func genericDecodeEncodeTest(t *testing.T, serializers ...polyio.Serializer) {
	for i, v := range serializers {
		r, w := io.Pipe()
		br := iotest.OneByteReader(r)
		go func() {
			if err := polyio.Encode(w, v); err != nil {
				t.Errorf("failed to encode %dth element (%T): %+v", i, v, err)
			}
			w.Close()
		}()

		dest := reflect.New(reflect.TypeOf(v).Elem())
		err := polyio.Decode(br, dest.Interface().(polyio.Serializer))
		r.Close()
		if err != nil {
			t.Errorf("failed to decode %dth element (%T): %+v", i, v, err)
		} else {
			_v := dest.Interface()
			assert.Equalf(t, v, _v, "comparing element %d", i)
		}
	}
}

// GenericBrokenPipeTest tests that encoding and decoding on broken streams fails.
func GenericBrokenPipeTest(t *testing.T, serializers ...polyio.Serializer) {
	for i, v := range serializers {
		r, w := io.Pipe()
		_ = w.Close()
		if err := v.Encode(w); err == nil {
			t.Errorf("encoding on closed writer should fail, but does not. %dth element (%T)", i, v)
		}

		_ = r.Close()
		if err := v.Decode(r); err == nil {
			t.Errorf("decoding on closed reader should fail, but does not. %dth element (%T)", i, v)
		}
	}
}
