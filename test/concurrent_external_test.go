// SPDX-License-Identifier: Apache-2.0

package test_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ctxtest "polycry.pt/poly-go/context/test"
	"polycry.pt/poly-go/test"
)

const timeout = 200 * time.Millisecond

func TestConcurrentT_Wait(t *testing.T) {
	t.Run("0 names", func(t *testing.T) {
		ct := test.NewConcurrent(t)
		require.Panics(t, func() { ct.Wait() })
	})

	t.Run("unknown name", func(t *testing.T) {
		ct := test.NewConcurrent(t)
		ct.Stage("known", func(t test.ConcT) {
		})
		ctxtest.AssertNotTerminates(t, timeout, func() { ct.Wait("unknown") })
	})

	t.Run("known name", func(t *testing.T) {
		ct := test.NewConcurrent(t)
		go ct.Stage("known", func(test.ConcT) {
			time.Sleep(timeout / 2)
		})
		ctxtest.AssertTerminates(t, timeout, func() { ct.Wait("known") })
	})

	t.Run("context expiry", func(t *testing.T) {
		ctxtest.AssertTerminates(t, timeout, func() {
			test.AssertFatal(t, func(t test.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				test.NewConcurrentCtx(ctx, t).Stage("", func(test.ConcT) {
					time.Sleep(timeout)
				})
			})
		})
	})
}

func TestConcurrentT_FailNow(t *testing.T) {
	t.Run("idempotence", func(t *testing.T) {
		var ct *test.ConcurrentT

		// Test that NewConcurrent.FailNow() calls T.FailNow().
		test.AssertFatal(t, func(t test.T) {
			ct = test.NewConcurrent(t)
			ct.FailNow()
		})

		// Test that after that, FailNow() calls runtime.Goexit().
		assert.True(t, test.CheckGoexit(ct.FailNow),
			"redundant FailNow() must call runtime.Goexit()")
	})

	t.Run("hammer", func(t *testing.T) {
		const parallel = 12
		for tries := 0; tries < 512; tries++ {
			test.AssertFatal(t, func(t test.T) {
				ct := test.NewConcurrent(t)
				for g := 0; g < parallel; g++ {
					go ct.StageN("concurrent", parallel, func(t test.ConcT) {
						t.FailNow()
					})
				}
				ct.Wait("concurrent")
			})
		}
	})
}

func TestConcurrentT_StageN(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		ct := test.NewConcurrent(t)
		var executed, returned sync.WaitGroup
		executed.Add(2)
		returned.Add(2)

		for i := 0; i < 2; i++ {
			go func() {
				ct.StageN("stage", 2, func(t test.ConcT) {
					executed.Done()
				})
				returned.Done()
			}()
		}

		ctxtest.AssertTerminates(t, timeout, executed.Wait)
		ctxtest.AssertTerminates(t, timeout, returned.Wait)
	})

	t.Run("n*m happy", func(t *testing.T) {
		N := 100
		M := 100

		ct := test.NewConcurrent(t)

		for g := 0; g < N; g++ {
			go func(g int) {
				for stage := 0; stage < M; stage++ {
					if g&1 == 0 {
						ct.StageN(strconv.Itoa(stage), N/2, func(t test.ConcT) {
						})
					} else {
						ct.Wait(strconv.Itoa(stage))
					}
				}
			}(g)
		}
	})

	t.Run("n*m sad", func(t *testing.T) {
		N := 100
		M := 100
		test.AssertFatal(t, func(t test.T) {
			ct := test.NewConcurrent(t)
			var wg sync.WaitGroup
			wg.Add(N)
			for g := 0; g < N; g++ {
				go func(g int) {
					defer wg.Done()
					for stage := 0; stage < M; stage++ {
						ct.StageN(strconv.Itoa(stage), N, func(t test.ConcT) {
							if g == N/2 {
								t.FailNow()
							}
						})
					}
				}(g)
			}

			wg.Wait()
		})
	})

	t.Run("too few goroutines", func(t *testing.T) {
		ct := test.NewConcurrent(t)
		ctxtest.AssertNotTerminates(t, timeout, func() {
			ct.StageN("stage", 2, func(test.ConcT) {})
		})
	})

	t.Run("too many goroutines", func(t *testing.T) {
		ct := test.NewConcurrent(t)
		go ct.StageN("stage", 2, func(test.ConcT) {})
		ct.StageN("stage", 2, func(test.ConcT) {})
		assert.Panics(t, func() {
			ct.StageN("stage", 2, func(test.ConcT) {})
		})
	})

	t.Run("inconsistent N", func(t *testing.T) {
		ct := test.NewConcurrent(t)
		var created sync.WaitGroup
		created.Add(1)

		go ct.StageN("stage", 2, func(test.ConcT) {
			created.Done()
		})

		created.Wait()
		assert.Panics(t, func() {
			ct.StageN("stage", 3, func(test.ConcT) {})
		})
	})

	t.Run("panic", func(t *testing.T) {
		test.AssertFatal(t, func(t test.T) {
			ct := test.NewConcurrent(t)
			ct.Stage("stage", func(test.ConcT) { panic(nil) })
		})
	})
}

func TestConcurrentT_Barrier(t *testing.T) {
	const N = 64
	const M = 32

	t.Run("happy", func(t *testing.T) {
		ct := test.NewConcurrent(t)

		for i := 0; i < N; i++ {
			go ct.StageN("loop", N, func(t test.ConcT) {
				for j := 0; j < M; j++ {
					t.BarrierN(fmt.Sprintf("barrier %d", j), N)
				}
			})
		}
		ct.Wait("loop")
	})

	t.Run("fail", func(t *testing.T) {
		test.AssertFatal(t, func(t test.T) {
			ct := test.NewConcurrent(t)

			for i := 0; i < N; i++ {
				i := i
				go ct.StageN("loop", N, func(t test.ConcT) {
					if i == N/2 {
						t.FailBarrierN("barrier", N)
					} else {
						t.BarrierN("barrier", N)
					}
				})
			}
			ct.Wait("loop")
		})
	})
}

func TestConcurrentT_PropagateBarrierFail(t *testing.T) {
	t.Run("FailBarrier causes other barrier to fail", func(t *testing.T) {
		ctxtest.AssertTerminatesQuickly(t, func() {
			test.AssertFatal(t, func(t test.T) {
				ct := test.NewConcurrent(t)
				go ct.FailBarrier("a")
				ct.BarrierN("b", 1000)
			})
		})
	})

	t.Run("Failed stage causes barrier to fail", func(t *testing.T) {
		ctxtest.AssertTerminatesQuickly(t, func() {
			test.AssertFatal(t, func(t test.T) {
				ct := test.NewConcurrent(t)
				go ct.Stage("a", func(t test.ConcT) {
					t.FailNow()
				})
				ct.BarrierN("b", 1000)
			})
		})
	})

	t.Run("FailNow causes Wait to fail", func(t *testing.T) {
		executed := false
		ctxtest.AssertTerminatesQuickly(t, func() {
			test.AssertFatal(t, func(t test.T) {
				ct := test.NewConcurrent(t)
				go ct.Stage("a", func(t test.ConcT) {
					t.FailNow()
				})
				ct.Wait("b")
				executed = true
			})
		})
		require.False(t, executed, "Wait() must call Goexit() on failure")
	})

	t.Run("FailNow must cause Barrier to fail", func(t *testing.T) {
		executed := false
		ctxtest.AssertTerminatesQuickly(t, func() {
			test.AssertFatal(t, func(t test.T) {
				ct := test.NewConcurrent(t)
				go ct.Stage("a", func(t test.ConcT) {
					t.BarrierN("c", 2)
					executed = true
				})
				go ct.Stage("b", func(t test.ConcT) {
					t.FailNow()
					t.BarrierN("c", 2)
				})
				ct.Wait("a", "b")
			})
		})
		require.False(t, executed, "Wait() must call Goexit() on failure")
	})
}
