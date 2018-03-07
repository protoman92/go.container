package gocollection

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func testListConcurrentOps(tb testing.TB, cl ConcurrentList) {
	keyCount := 1000
	keys := make([]interface{}, keyCount)

	for ix := range keys {
		keys[ix] = strconv.Itoa(ix)
	}

	sleepRandomizer := func() time.Duration {
		var min, max time.Duration = 1e8, 5e8
		duration := min + time.Duration(rand.Int63n(int64(max-min)))
		return duration
	}

	params := &ConcurrentOpsParams{
		concurrentCollection: cl,
		log:                  false,
		keys:                 keys,
		opSleepDuration:      sleepRandomizer,
	}

	wgAccess := setupConcurrentCollectionOps(params)

	// Get
	for ix := range keys {
		wgAccess().Add(1)

		go func(ix int) {
			time.Sleep(sleepRandomizer())

			cl.GetAtAsync(ix, func(e interface{}, found bool) {
				if params.log {
					fmt.Printf("Got element %v, found: %t\n", e, found)
				}

				wgAccess().Done()
			})
		}(ix)
	}

	// Modify
	for ix := range keys {
		wgAccess().Add(1)

		go func(ix int) {
			time.Sleep(sleepRandomizer())

			cl.RemoveAtAsync(ix, func(e interface{}, found bool) {
				if params.log {
					fmt.Printf("Deleted element %v, found: %t\n", e, found)
				}

				wgAccess().Done()
			})
		}(ix)
	}

	for ix := range keys {
		wgAccess().Add(1)

		indexes := make([]int, 0)

		for jx := range keys {
			indexes = append(indexes, jx)
		}

		go func(ix int) {
			time.Sleep(sleepRandomizer())

			cl.RemoveAllAtAsync(func(deleted int) {
				if params.log {
					fmt.Printf("Deleted %d elements\n", deleted)
				}

				wgAccess().Done()
			}, indexes...)
		}(ix)
	}

	for ix := range keys {
		wgAccess().Add(1)

		go func(ix int) {
			time.Sleep(sleepRandomizer())

			cl.SetAtAsync(ix, keys[ix], func(e interface{}, found bool) {
				if params.log {
					fmt.Printf("Prev element %v, found: %t\n", e, found)
				}

				wgAccess().Done()
			})
		}(ix)
	}

	wgAccess().Wait()
	fmt.Printf("Final storage state: %v", cl)
}

func benchmarkListConcurrentOps(b *testing.B, clFn func() ConcurrentList) {
	for i := 0; i < b.N; i++ {
		testListConcurrentOps(b, clFn())
	}
}

func BenchmarkLockSliceListConcurrentOps(b *testing.B) {
	benchmarkListConcurrentOps(b, func() ConcurrentList {
		sl := NewDefaultSliceList()
		return NewLockConcurrentList(sl)
	})
}

func TestLockSliceListConcurrentOps(t *testing.T) {
	sl := NewDefaultSliceList()
	cl := NewLockConcurrentList(sl)
	testListConcurrentOps(t, cl)
}
