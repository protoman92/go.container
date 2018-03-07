package gomap

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

type ConcurrentMapOpsParams struct {
	concurrentMap   ConcurrentMap
	log             bool
	keyCount        int
	opSleepDuration func() time.Duration
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
			time.Sleep(params.opSleepDuration())

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
			time.Sleep(params.opSleepDuration())

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
			time.Sleep(params.opSleepDuration())

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
			time.Sleep(params.opSleepDuration())

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
			time.Sleep(params.opSleepDuration())

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
			time.Sleep(params.opSleepDuration())

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
	fmt.Printf("Final map %v\n", cm)
}

func testConcurrentMapConcurrentOps(tb testing.TB, cm ConcurrentMap) {
	sleepRandomizer := func() time.Duration {
		var min, max time.Duration = 1e8, 5e8
		duration := min + time.Duration(rand.Int63n(int64(max-min)))
		return duration
	}

	params := &ConcurrentMapOpsParams{
		concurrentMap:   cm,
		log:             false,
		keyCount:        1000,
		opSleepDuration: sleepRandomizer,
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
	bm := NewDefaultBasicMap()
	cm := NewChannelConcurrentMap(bm)
	testConcurrentMapConcurrentOps(t, cm)
}

func TestLockConcurrentMapConcurrentOps(t *testing.T) {
	bm := NewDefaultBasicMap()
	cm := NewLockConcurrentMap(bm)
	testConcurrentMapConcurrentOps(t, cm)
}
