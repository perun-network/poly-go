// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentT_Wait(t *testing.T) {
	t.Run("failed stage", func(t *testing.T) {
		AssertFatal(t, func(t T) {
			ct := NewConcurrent(t)
			s := ct.spawnStage("stage", 1)
			s.failed.Set()
			s.pass()

			ct.Wait("stage")
		})
	})
}

func TestStage_FailNow(t *testing.T) {
	t.Run("first fail", func(t *testing.T) {
		AssertFatal(t, func(t T) {
			ct := NewConcurrent(t)
			s := ct.spawnStage("stage", 1)
			s.FailNow()
		})
	})

	t.Run("second fail", func(t *testing.T) {
		ct := NewConcurrent(nil)
		ct.failed = true
		s := ct.spawnStage("stage", 1)
		assert.True(t, CheckGoexit(s.FailNow))
	})
}
