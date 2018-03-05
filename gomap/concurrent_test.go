package gomap

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func Test_ExecutingConcurrentOpsOnConcurrentMap_ShouldWork(t *testing.T) {
	/// Setup
	cm := NewDefaultBasicConcurrentMap()
	keys := make([]string, 1000)

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
				cm.GetAsync(key, func(value interface{}, found bool) {
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

	accessWaitGroup().Wait()

	/// Then
	// It does not matter what the values here are - running the tests with race
	// mode will automatically fail if concurrent operations are not performed
	// correctly.
	fmt.Printf("Final length: %d\n", cm.Length())
	fmt.Printf("Final map: %v\n", cm.UnderlyingMap())
	fmt.Printf("Final storage: %v\n", cm.UnderlyingStorage())
}
