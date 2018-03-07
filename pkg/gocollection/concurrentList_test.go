package gocollection

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func testListConcurrentOps(tb testing.TB, cl List) {
	keyCount := 500
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
	go func() {
		for ix := range keys {
			wgAccess().Add(1)

			go func(ix int) {
				time.Sleep(sleepRandomizer())

				if e, found := cl.GetAt(ix); params.log {
					fmt.Printf("Got element %v, found: %t\n", e, found)
				}

				wgAccess().Done()
			}(ix)
		}
	}()

	go func() {
		for ix := range keys {
			wgAccess().Add(1)

			go func(key interface{}) {
				time.Sleep(sleepRandomizer())

				if ix, found := cl.IndexOf(key); params.log {
					fmt.Printf("Got index %d, found: %t\n", ix, found)
				}

				wgAccess().Done()
			}(keys[ix])
		}
	}()

	// Modify
	go func() {
		for ix := range keys {
			wgAccess().Add(1)

			go func(ix int) {
				time.Sleep(sleepRandomizer())

				if e, found := cl.RemoveAt(ix); params.log {
					fmt.Printf("Deleted element %v, found: %t\n", e, found)
				}

				wgAccess().Done()
			}(ix)
		}
	}()

	go func() {
		for ix := range keys {
			wgAccess().Add(1)

			indexes := make([]int, 0)

			for jx := range keys {
				indexes = append(indexes, jx)
			}

			go func(ix int) {
				time.Sleep(sleepRandomizer())

				if deleted := cl.RemoveAllAt(indexes...); params.log {
					fmt.Printf("Deleted %d elements\n", deleted)
				}

				wgAccess().Done()
			}(ix)
		}
	}()

	go func() {
		for ix := range keys {
			wgAccess().Add(1)

			go func(ix int) {
				time.Sleep(sleepRandomizer())

				if e, found := cl.SetAt(ix, keys[ix]); params.log {
					fmt.Printf("Prev element %v, found: %t\n", e, found)
				}

				wgAccess().Done()
			}(ix)
		}
	}()

	time.Sleep(time.Millisecond)
	wgAccess().Wait()
	fmt.Printf("Final storage state: %v", cl)
}

func benchmarkListConcurrentOps(b *testing.B, clFn func() List) {
	for i := 0; i < b.N; i++ {
		testListConcurrentOps(b, clFn())
	}
}

func BenchmarkLockSliceListConcurrentOps(b *testing.B) {
	benchmarkListConcurrentOps(b, func() List {
		sl := NewDefaultSliceList()
		return NewLockConcurrentList(sl)
	})
}

func TestLockSliceListConcurrentOps(t *testing.T) {
	sl := NewDefaultSliceList()
	cl := NewLockConcurrentList(sl)
	testListConcurrentOps(t, cl)
}
