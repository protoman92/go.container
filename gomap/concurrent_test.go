package gomap

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func testConcurrentMapBasicOps(tb testing.TB, cm ConcurrentMap) {
	/// Setup
	key := "Key"
	value := "Value"

	/// When & Then
	cm.Set(key, value)

	if !cm.Contains(key) {
		tb.Errorf("Should contain %v", key)
	}

	value1, found := cm.Get(key)

	if value1 != value || !found {
		tb.Errorf("Should contain %v with value %v", key, value)
	}

	cm.Delete(key)
	_, found1 := cm.Get(key)

	if found1 && cm.Contains(key) {
		tb.Errorf("Should not contain %v", key)
	}

	cm.Set(key, value)
	length := cm.Length()

	if length != 1 {
		tb.Errorf("Should have length 1, but got %d", length)
	}

	cm.Clear()

	if !cm.IsEmpty() {
		tb.Errorf("Should be empty")
	}

	fmt.Printf("Final storage: %v\n", cm.UnderlyingStorage())
}

func TestChannelConcurrentMapBasicOps(t *testing.T) {
	t.Parallel()
	bm := NewDefaultBasicMap()
	cm := NewChannelConcurrentMap(bm)
	testConcurrentMapBasicOps(t, cm)
}

func TestLockConcurrentMapBasicOps(t *testing.T) {
	t.Parallel()
	bm := NewDefaultBasicMap()
	cm := NewLockConcurrentMap(bm)
	testConcurrentMapBasicOps(t, cm)
}

func testConcurrentMapConcurrentOps(tb testing.TB, cm ConcurrentMap) {
	/// Setup
	keys := make([]string, 200)

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
	go func() {
		for _, key := range keys {
			accessWaitGroup().Add(1)

			go func(key string) {
				cm.GetAsync(key, func(value Value, found bool) {
					accessWaitGroup().Done()
				})
			}(key)
		}
	}()

	go func() {
		for _, key := range keys {
			accessWaitGroup().Add(1)

			go func(key string) {
				cm.DeleteAsync(key, func(len int) {
					accessWaitGroup().Done()
				})
			}(key)
		}
	}()

	go func() {
		for i := 0; i < len(keys); i++ {
			accessWaitGroup().Add(1)

			go func() {
				cm.ClearAsync(func() {
					accessWaitGroup().Done()
				})
			}()
		}
	}()

	// Get
	go func() {
		for _, key := range keys {
			accessWaitGroup().Add(1)

			go func(key string) {
				cm.SetAsync(key, key, func(len int) {
					accessWaitGroup().Done()
				})
			}(key)
		}
	}()

	go func() {
		for _, key := range keys {
			accessWaitGroup().Add(1)

			go func(key string) {
				cm.ContainsAsync(key, func(found bool) {
					accessWaitGroup().Done()
				})
			}(key)
		}
	}()

	go func() {
		for i := 0; i < len(keys); i++ {
			accessWaitGroup().Add(1)

			go func() {
				cm.LengthAsync(func(len int) {
					accessWaitGroup().Done()
				})
			}()
		}
	}()

	go func() {
		for i := 0; i < len(keys); i++ {
			accessWaitGroup().Add(1)

			go func() {
				cm.IsEmptyAsync(func(empty bool) {
					accessWaitGroup().Done()
				})
			}()
		}
	}()

	// Access
	go func() {
		for i := 0; i < len(keys); i++ {
			accessWaitGroup().Add(1)

			go func() {
				cm.UnderlyingStorageAsync(func(storage map[Key]Value) {
					accessWaitGroup().Done()
				})
			}()
		}
	}()

	accessWaitGroup().Wait()

	/// Then
	// It does not matter what we assert here - running the tests with race mode
	// will automatically fail if concurrent ops are not performed correctly.
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
