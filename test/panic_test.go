// SPDX-License-Identifier: Apache-2.0

package test

import "testing"

func TestCheckPanic(t *testing.T) {
	// Test whether panic calls are properly detected and whether the supplied
	// value is also properly recorded.
	if p, v := CheckPanic(func() { panic("panicvalue") }); !p || v != "panicvalue" {
		t.Error("Failed to detect panic!")
	}

	// Test whether panic(nil) calls are detected.
	if p, v := CheckPanic(func() { panic(nil) }); !p || v != nil {
		t.Error("Failed to detect panic(nil)!")
	}

	// Test whether the absence of panic calls is properly detected.
	if p, v := CheckPanic(func() {}); p || v != nil {
		t.Error("False positive panic detection!")
	}
}
