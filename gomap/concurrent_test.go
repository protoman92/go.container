package gomap

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

type ConcurrentMapOpsParams struct {
	concurrentMap ConcurrentMap
	log           bool
	keyCount      int
}

func setupConcurrentMapOps(params *ConcurrentMapOpsParams) {
	/// Setup
	cm := params.concurrentMap
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

	/// When
	// Modify
	for _, key := range keys {
		accessWaitGroup().Add(1)

		go func(key string) {
			cm.SetAsync(key, key, func(prev Value, found bool) {
				if params.log {
					fmt.Printf("Set key %v-value %v. Prev value: %v\n", key, key, prev)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for _, key := range keys {
		accessWaitGroup().Add(1)

		go func(key string) {
			cm.DeleteAsync(key, func(found bool) {
				if params.log {
					fmt.Printf("Deleted key %v, found: %t\n", key, found)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for i := 0; i < len(keys); i++ {
		accessWaitGroup().Add(1)

		go func() {
			cm.ClearAsync(func() {
				if params.log {
					fmt.Printf("Cleared all contents\n")
				}

				accessWaitGroup().Done()
			})
		}()
	}

	// Get
	for _, key := range keys {
		accessWaitGroup().Add(1)

		go func(key string) {
			cm.GetAsync(key, func(value Value, found bool) {
				if params.log {
					fmt.Printf("Got %v for key %v, found: %t\n", value, key, found)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for _, key := range keys {
		accessWaitGroup().Add(1)

		go func(key string) {
			cm.ContainsAsync(key, func(found bool) {
				if params.log {
					fmt.Printf("Contains key %v: %t\n", key, found)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for i := 0; i < len(keys); i++ {
		accessWaitGroup().Add(1)

		go func() {
			cm.LengthAsync(func(len int) {
				if params.log {
					fmt.Printf("Current length: %d\n", len)
				}

				accessWaitGroup().Done()
			})
		}()
	}

	accessWaitGroup().Wait()

	/// Then
	// It does not matter what we assert here - running the tests with race mode
	// will automatically fail if concurrent ops are not performed correctly.
}

func testConcurrentMapConcurrentOps(tb testing.TB, cm ConcurrentMap) {
	params := &ConcurrentMapOpsParams{
		concurrentMap: cm,
		log:           false,
		keyCount:      500,
	}

	setupConcurrentMapOps(params)
}

func benchmarkConcurrentMapConcurrentOps(b *testing.B, cmFn func() ConcurrentMap) {
	for i := 0; i < b.N; i++ {
		testConcurrentMapConcurrentOps(b, cmFn())
	}
}

func BenchmarkChannelConcurrentMapConcurrentOps(b *testing.B) {
	benchmarkConcurrentMapConcurrentOps(b, func() ConcurrentMap {
		bm := NewDefaultBasicMap()
		return NewChannelConcurrentMap(bm)
	})
}

func BenchmarkLockConcurrentMapConcurrentOps(b *testing.B) {
	benchmarkConcurrentMapConcurrentOps(b, func() ConcurrentMap {
		bm := NewDefaultBasicMap()
		return NewLockConcurrentMap(bm)
	})
}

func TestChannelConcurrentMapConcurrentOps(t *testing.T) {
	t.Parallel()
	bm := NewDefaultBasicMap()
	cm := NewChannelConcurrentMap(bm)
	testConcurrentMapConcurrentOps(t, cm)
}

func TestLockConcurrentMapConcurrentOps(t *testing.T) {
	t.Parallel()
	bm := NewDefaultBasicMap()
	cm := NewLockConcurrentMap(bm)
	testConcurrentMapConcurrentOps(t, cm)
}
