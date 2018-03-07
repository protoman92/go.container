package gocollection

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type ConcurrentOpsParams struct {
	concurrentCollection ConcurrentCollection
	log                  bool
	keys                 []string
	opSleepDuration      func() time.Duration
}

func setupConcurrentCollectionOps(params *ConcurrentOpsParams) func() *sync.WaitGroup {
	/// Setup
	cl := params.concurrentCollection
	keys := params.keys

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
			time.Sleep(params.opSleepDuration())

			cl.AddAsync(key, func(added int) {
				if params.log {
					fmt.Printf("Added %v, added count: %v\n", key, added)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for _, key := range keys {
		accessWaitGroup().Add(1)

		go func(key string) {
			time.Sleep(params.opSleepDuration())

			cl.RemoveAsync(key, func(removed int) {
				if params.log {
					fmt.Printf("Deleted %v, removed count: %d\n", key, removed)
				}

				accessWaitGroup().Done()
			})
		}(key)
	}

	for i := 0; i < len(keys); i++ {
		accessWaitGroup().Add(1)

		go func() {
			time.Sleep(params.opSleepDuration())

			cl.ClearAsync(func() {
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
			time.Sleep(params.opSleepDuration())

			cl.LengthAsync(func(len int) {
				if params.log {
					fmt.Printf("Current length: %d\n", len)
				}

				accessWaitGroup().Done()
			})
		}()
	}

	return accessWaitGroup
}
