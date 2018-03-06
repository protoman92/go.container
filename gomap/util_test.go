package gomap

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type ConcurrentMapOpsParams struct {
	concurrentMap     ConcurrentMap
	log               bool
	keyCount          int
	setupWaitDuration time.Duration
}

func SetupConcurrentMapOps(params *ConcurrentMapOpsParams) {
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
	go func() {
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
	}()

	go func() {
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
	}()

	go func() {
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
	}()

	// Get
	go func() {
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
	}()

	go func() {
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
	}()

	go func() {
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
	}()

	go func() {
		for i := 0; i < len(keys); i++ {
			accessWaitGroup().Add(1)

			go func() {
				cm.IsEmptyAsync(func(empty bool) {
					if params.log {
						fmt.Printf("Is empty: %t\n", empty)
					}

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
					if params.log {
						fmt.Printf("Current storage %v\n", storage)
					}

					accessWaitGroup().Done()
				})
			}()
		}
	}()

	// Sleep for a short while to let all goroutines activate.
	time.Sleep(params.setupWaitDuration)
	accessWaitGroup().Wait()

	/// Then
	// It does not matter what we assert here - running the tests with race mode
	// will automatically fail if concurrent ops are not performed correctly.
}
