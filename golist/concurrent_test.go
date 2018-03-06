package golist

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

type ConcurrentListOpsParams struct {
	concurrentList ConcurrentList
	log            bool
	keyCount       int
}

func setupConcurrentListOps(params *ConcurrentListOpsParams) {
	/// Setup
	cl := params.concurrentList
	keys := make([]string, params.keyCount)

	for ix := range keys {
		keys[ix] = strconv.Itoa(ix)
	}

	waitGroup := &sync.WaitGroup{}
	mutex := sync.RWMutex{}

	accessWaitGroup := func() *sync.WaitGroup {
		mutex.Lock()
		defer mutex.Unlock()
		return waitGroup
	}

	/// When & Then
	// Modify
	for _, key := range keys {
		accessWaitGroup().Add(1)

		go func(key string) {
			cl.AddAsync(key, func(added int) {
				if params.log {
					fmt.Printf("Added %v. Total added: %v\n", key, added)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for _, key := range keys {
		accessWaitGroup().Add(1)

		go func(key string) {
			cl.RemoveAsync(key, func(found bool) {
				if params.log {
					fmt.Printf("Deleted %v, found: %t\n", key, found)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for i := 0; i < len(keys); i++ {
		accessWaitGroup().Add(1)

		go func() {
			cl.ClearAsync(func() {
				if params.log {
					fmt.Printf("Cleared all contents\n")
				}

				accessWaitGroup().Done()
			})
		}()
	}

	// Get
	for ix := range keys {
		accessWaitGroup().Add(1)

		go func(ix int) {
			cl.GetAsync(ix, func(e Element, found bool) {
				if params.log {
					fmt.Printf("Got %v for index %d, found: %t\n", e, ix, found)
				}

				accessWaitGroup().Done()
			})
		}(ix)
	}

	for _, key := range keys {
		accessWaitGroup().Add(1)

		go func(key string) {
			cl.ContainsAsync(key, func(found bool) {
				if params.log {
					fmt.Printf("Contains element %v: %t\n", key, found)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for i := 0; i < len(keys); i++ {
		accessWaitGroup().Add(1)

		go func() {
			cl.LengthAsync(func(len int) {
				if params.log {
					fmt.Printf("Current length: %d\n", len)
				}

				accessWaitGroup().Done()
			})
		}()
	}

	accessWaitGroup().Wait()
}

func testConcurrentListConcurrentOps(tb testing.TB, cl ConcurrentList) {
	params := &ConcurrentListOpsParams{
		concurrentList: cl,
		log:            false,
		keyCount:       500,
	}

	setupConcurrentListOps(params)
}

func benchmarkConcurrentListConcurrentOps(b *testing.B, clFn func() ConcurrentList) {
	for i := 0; i < b.N; i++ {
		testConcurrentListConcurrentOps(b, clFn())
	}
}

func BenchmarkLockConcurrentListConcurrentOps(b *testing.B) {
	benchmarkConcurrentListConcurrentOps(b, func() ConcurrentList {
		sl := NewSliceList()
		return NewLockConcurrentList(sl)
	})
}

func TestLockConcurrentListConcurrentOps(t *testing.T) {
	t.Parallel()
	sl := NewSliceList()
	cl := NewLockConcurrentList(sl)
	testConcurrentListConcurrentOps(t, cl)
}
