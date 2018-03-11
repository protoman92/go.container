package gocollection

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func testListConcurrentOps(tb testing.TB, cl List) {
	keyCount := 200
	keys := make([]interface{}, keyCount)

	for ix := range keys {
		keys[ix] = strconv.Itoa(ix)
	}

	sleepRandomizer := func() time.Duration {
		return time.Duration(1e6)
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

				if e, found := cl.GetFirst(); params.log {
					fmt.Printf("Got first element %v, found: %t\n", e, found)
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

				if ix, e, found := cl.GetFirstFunc(func(ix int, e interface{}) bool {
					return e == key
				}); params.log {
					fmt.Printf("Got first %v at index %d, found: %t\n", e, ix, found)
				}

				wgAccess().Done()
			}(keys[ix])
		}
	}()

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

	go func() {
		for ix := range keys {
			wgAccess().Add(1)

			go func(key interface{}) {
				time.Sleep(sleepRandomizer())

				if ix, e, found := cl.IndexOfFunc(func(ix int, e interface{}) bool {
					return e == key
				}); params.log {
					fmt.Printf("Got index %d with element %v, found: %t\n", ix, e, found)
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
	fmt.Printf("Final storage state: %v\n", cl)
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
