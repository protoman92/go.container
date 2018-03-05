package gomap

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func Test_ExecutingBasicOpsOnConcurrentMap_ShouldWork(t *testing.T) {
	/// Setup
	cm := NewDefaultBasicConcurrentMap()
	key := "Key"
	value := "Value"

	/// When & Then
	cm.Set(key, value)

	if !cm.Contains(key) {
		t.Errorf("Should contain %v", key)
	}

	value1, found := cm.Get(key)

	if value1 != value || !found {
		t.Errorf("Should contain %v with value %v", key, value)
	}

	cm.Delete(key)
	_, found1 := cm.Get(key)

	if found1 && cm.Contains(key) {
		t.Errorf("Should not contain %v", key)
	}

	cm.Set(key, value)
	length := cm.Length()

	if length != 1 {
		t.Errorf("Should have length 1, but got %d", length)
	}

	cm.Clear()

	if !cm.IsEmpty() {
		t.Errorf("Should be empty")
	}
}

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
				cm.ContainsAsync(key, func(found bool) {
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

	accessWaitGroup().Wait()

	/// Then
	// It does not matter what the values here are - running the tests with race
	// mode will automatically fail if concurrent operations are not performed
	// correctly.
	fmt.Printf("Final length: %d\n", cm.Length())
	fmt.Printf("Final map: %v\n", cm.UnderlyingMap())
	fmt.Printf("Final storage: %v\n", cm.UnderlyingStorage())
}