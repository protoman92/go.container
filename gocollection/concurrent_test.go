package gocollection

import (
	"fmt"
	"strconv"
	"sync"
)

type ConcurrentCollectionOpsParams struct {
	concurrentCollection ConcurrentCollection
	log                  bool
	keyCount             int
}

func setupConcurrentCollectionOps(params *ConcurrentCollectionOpsParams) {
	/// Setup
	cl := params.concurrentCollection
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
